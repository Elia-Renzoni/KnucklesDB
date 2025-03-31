package swim

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"knucklesdb/wal"
	"net"
	"time"
)

const MAX_GOSSIP_ATTEMPS int = 3
const RUNE_MEMBERSHIP_LIST_SEPARATOR rune = ','

const (
	SPREAD_MEMBERSHIP int = iota * 1
	SPREAD_UPDATES
)

type Dissemination struct {
	conn                      net.Conn
	logger                    *wal.InfoLogger
	errorLogger               *wal.ErrorsLogger
	gossipGlobalContext       context.Context
	timeoutTime               time.Duration
	cluster                   *ClusterManager
	marshaler                 *ProtocolMarshaer
	ack                       AckMessage
	gossipQuorum              int
	gossipQuorumSpreadingList int
}

func NewDissemination(timeoutTime time.Duration, logger *wal.InfoLogger, errorLogger *wal.ErrorsLogger, cluster *ClusterManager,
	marshaler *ProtocolMarshaer) *Dissemination {
	return &Dissemination{
		logger:              logger,
		errorLogger:         errorLogger,
		gossipGlobalContext: context.Background(),
		timeoutTime:         timeoutTime,
		cluster:             cluster,
		marshaler:           marshaler,
		gossipQuorum:        0,
	}
}

func (d *Dissemination) SpreadMembershipList(membershipList []*Node, fanoutList []string) {
	for attemp := 0; attemp < MAX_GOSSIP_ATTEMPS; attemp++ {
		for index := range fanoutList {
			encodeClusterMetadata, err := d.marshalMembershipList(membershipList)
			d.errorLogger.ReportError(err)
			d.send(fanoutList[index], encodeClusterMetadata, SPREAD_MEMBERSHIP)
		}

		// check if the majority quorum is reached.
		if d.gossipQuorumSpreadingList >= (len(fanoutList)/2)+1 {
			break
		}
	}

	// reset the qourum counter
	d.gossipQuorumSpreadingList = 0
}

func (d *Dissemination) TransformMembershipList(cluster MembershipListMessage) []*Node {
	var (
		clusterNodes []*Node = make([]*Node, 0)
		node         MembershipEntry
	)

	for _, node := range cluster.List {
		convertedNode := NewNode(node.NodeAddress, node.NodeListenPort, node.NodeStatus)
		clusterNodes = append(clusterNodes, convertedNode)
	}

	return clusterNodes
}

func (d *Dissemination) IsMembershipListDifferent(receivedMembershipList []*Node) bool {
	var (
		different bool = true
	)

	for remoteNodeIndex := range receivedMembershipList {
		for localNodeIndex := range d.cluster.clusterMetadata {
			switch {
			case receivedMembershipList[remoteNodeIndex].nodeAddress != d.cluster.clusterMetadata[localNodeIndex].nodeAddress:
				fallthrough
			case receivedMembershipList[remoteNodeIndex].nodeListenPort != d.cluster.clusterMetadata[localNodeIndex].nodeListenPort:
				fallthrough
			case receivedMembershipList[remoteNodeIndex].nodeStatus != d.cluster.clusterMetadata[localNodeIndex].nodeStatus:
				different = false
			default:
				different = true
			}
		}
	}

	return different
}

func (d *Dissemination) MergeMembershipList(clusterMetadata []*Node) {
	var differencies, length = d.getDifferencies(clusterMetadata)

	if length != 0 {
		for _, nodeToJoin := range differencies {
			d.cluster.clusterMetadata = append(d.cluster.clusterMetadata, nodeToJoin)
		}
	}
}

func (d *Dissemination) getDifferencies(receivedClusterMembers []*Node) ([]*Node, int) {
	var (
		diffSlice []*Node = make([]*Node, 0)
		different bool
	)

	for _, receivedNode := range receivedClusterMembers {
		for _, node := range d.cluster.clusterMetadata {
			switch {
			case receivedNode.nodeAddress == node.nodeAddress:
				fallthrough
			case receivedNode.nodeListenPort == node.nodeListenPort:
				fallthrough
			case receivedNode.nodeStatus == node.nodeStatus:
				different = false
			default:
				different = true
			}
		}

		if different {
			newNode := NewNode(receivedNode.nodeAddress, receivedNode.nodeListenPort, receivedNode.nodeStatus)
			diffSlice = append(diffSlice, newNode)
		}
	}

	return diffSlice, len(diffSlice)
}

func (d *Dissemination) SpreadMembershipListUpdates(fanoutList []string, updateToSpread *Node) {
	for i := 0; i < MAX_GOSSIP_ATTEMPS; i++ {
		for index := range fanoutList {
			encodedUpdate, _ := d.marshaler.MarshalSingleNodeUpdate(updateToSpread.nodeAddress, updateToSpread.nodeListenPort, updateToSpread.nodeStatus)
			d.send(fanoutList[index], encodedUpdate, SPREAD_UPDATES)
		}

		// check if the majority quorum is reached.
		if d.gossipQuorum >= (len(fanoutList)/2)+1 {
			break
		}
	}

	// reset the quorum counter
	d.gossipQuorum = 0
}

func (d *Dissemination) send(nodeAddress string, gossipMessage []byte, operation int) {
	ctx, cancel := context.WithTimeout(d.gossipGlobalContext, d.timeoutTime)
	defer cancel()
	conn, err := net.Dial("tcp", nodeAddress)
	if err != nil {
		d.errorLogger.ReportError(err)
		return
	}
	defer conn.Close()

	conn.Write(gossipMessage)

	data := make([]byte, 2024)

	select {
	case <-ctx.Done():
		d.errorLogger.ReportError(errors.New("Gossip Send Failed due to Context Timeout"))
	default:
		count, _ := conn.Read(data)
		json.Unmarshal(data[:count], &d.ack)
		switch operation {
		case SPREAD_MEMBERSHIP:
			if d.ack.AckContent == 1 {
				d.gossipQuorumSpreadingList += 1
			}
		case SPREAD_UPDATES:
			if d.ack.AckContent == 1 {
				d.gossipQuorum += 1
			}
		}
	}
}

func (d *Dissemination) marshalMembershipList(clusterData []*Node) ([]byte, error) {
	var (
		encodedMembershipList bytes.Buffer
		nodesObject           bytes.Buffer
		entry, header, tail   []byte
		err                   error
	)

	header, err = json.Marshal(map[string]string{
		"type": "membership",
	})
	encodedMembershipList.Write(header)

	for index := range clusterData {
		entry, err = json.Marshal(map[string]any{
			"address": clusterData[index].nodeAddress,
			"port":    clusterData[index].nodeListenPort,
			"status":  clusterData[index].nodeStatus,
		})

		if err != nil {
			return nil, err
		}

		nodesObject.Write(entry)
		nodesObject.WriteRune(RUNE_MEMBERSHIP_LIST_SEPARATOR)
	}

	tail, err = json.Marshal(map[string][]byte{
		"list": nodesObject.Bytes(),
	})

	encodedMembershipList.Write(tail)

	return encodedMembershipList.Bytes(), err
}
