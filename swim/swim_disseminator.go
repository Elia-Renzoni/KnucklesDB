package swim

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"knucklesdb/wal"
	"net"
	"time"
	"strconv"
	_"fmt"
	"sync"
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
	cluster                   *Cluster
	marshaler                 *ProtocolMarshaer
	ack                       AckMessage
	mutex  *sync.Mutex
	gossipQuorum              int
	gossipQuorumSpreadingList int
}

func NewDissemination(timeoutTime time.Duration, logger *wal.InfoLogger, errorLogger *wal.ErrorsLogger, cluster *Cluster,
	marshaler *ProtocolMarshaer, mutex *sync.Mutex) *Dissemination {
	return &Dissemination{
		logger:              logger,
		errorLogger:         errorLogger,
		gossipGlobalContext: context.Background(),
		timeoutTime:         timeoutTime,
		cluster:             cluster,
		marshaler:           marshaler,
		mutex: mutex,
		gossipQuorum:        0,
	}
}

func (d *Dissemination) SpreadMembershipList(membershipList []*Node, fanoutList []string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	encodeClusterMetadata, err := d.marshalMembershipList(membershipList)
	if err != nil {
		d.errorLogger.ReportError(err)
		return
	}

	for index := range fanoutList {
		d.send(fanoutList[index], encodeClusterMetadata, SPREAD_MEMBERSHIP)
	}
}

func (d *Dissemination) TransformMembershipList(cluster MembershipListMessage) []*Node {
	var (
		clusterNodes []*Node = make([]*Node, 0)
	)

	for _, node := range cluster.List {
		port, _ := strconv.Atoi(node.NodeListenPort)
		convertedNode := NewNode(node.NodeAddress, port, node.NodeStatus)
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

func (d *Dissemination) IsUpdateDifferent(update *Node) bool {
	var (
		different bool = false
	)

	for _, node := range d.cluster.clusterMetadata {
		if node.nodeAddress == update.nodeAddress {
			if node.nodeListenPort == update.nodeListenPort {
				if node.nodeStatus != update.nodeStatus {
					different = true
				}
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

func (d *Dissemination) MergeUpdates(update *Node) {
	for _, node := range d.cluster.clusterMetadata {
		switch {
		case node.nodeAddress == update.nodeAddress:
			fallthrough
		case node.nodeListenPort == update.nodeListenPort:
			node.nodeStatus = update.nodeStatus
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
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for index := range fanoutList {
		encodedUpdate, _ := d.marshaler.MarshalSingleNodeUpdate(updateToSpread.nodeAddress, updateToSpread.nodeListenPort, updateToSpread.nodeStatus)
		d.send(fanoutList[index], encodedUpdate, SPREAD_UPDATES)
	}
}

func (d *Dissemination) send(nodeAddress string, gossipMessage []byte, operation int) {
	ctx, cancel := context.WithTimeout(d.gossipGlobalContext, d.timeoutTime)
	defer cancel()

	var (
		err error
	)

	d.conn, err = net.Dial("tcp", nodeAddress)
	if err != nil {
		d.errorLogger.ReportError(err)
		return
	}
	defer d.conn.Close()

	d.conn.Write(gossipMessage)

	data := make([]byte, 2024)

	select {
	case <-ctx.Done():
		d.errorLogger.ReportError(errors.New("Gossip Send Failed due to Context Timeout"))
	default:
		count, _ := d.conn.Read(data)
		json.Unmarshal(data[:count], &d.ack)

		d.logger.ReportInfo("Ack Message Arrived from Neighbours")

		switch operation {
		case SPREAD_MEMBERSHIP:
			switch d.ack.AckContent {
			case 1, 0:
				d.gossipQuorumSpreadingList += 1
			}
		case SPREAD_UPDATES:
			switch d.ack.AckContent {
			case 1, 0:
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
