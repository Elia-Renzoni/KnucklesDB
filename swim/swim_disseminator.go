package swim

import (
	"net"
	"knucklesdb/wal"
	"context"
	"time"
	"encoding/json"
)

const MAX_GOSSIP_ATTEMPS int = 3

type Dissemination struct {
	conn net.Conn
	clusterNodes []MembershipEntry
	logger *wal.InfoLogger
	errorLogger *wal.ErrorsLogger
	gossipGlobalContext context.Context 
	timeoutTime time.Duration
	cluster *ClusterManager
	marshaler *ProtocolMarshaer
	ack       AckMessage
	gossipQuorum int
	gossipQuorumSpreadingList int 
}

type MembershipEntry struct {
	NodeAddress string `json:"address"`
	NodeListenPort string `json:"port"`
	NodeStatus int `json:"status"`
}

func NewDissemination(timeoutTime time.Duration, logger *wal.InfoLogger, errorLogger *wal.ErrorsLogger, cluster *ClusterManager,
	marshaler *ProtocolMarshaer) *Dissemination {
	return &Dissemination{
		logger: logger, 
		errorLogger: errorLogger,
		gossipGlobalContext: context.Background(),
		timeoutTime: timeoutTime,
		cluster: cluster,
		marshaler: marshaler,
		gossipQuorum: 0,
	}
}

func (d *Dissemination) SpreadMembershipList(membershipList []*Node, fanoutList []string) {
	for attemp := 0; attemp < MAX_GOSSIP_ATTEMPS; attemp++ {
		for index := range fanoutList {
			encodeClusterMetadata, err := d.marshalMembershipList(membershipList)
			d.send(fanoutList[index], encodeClusterMetadata)
		}

		// TODO: Add quorum policies
	}
}

func (d *Dissemination) IsMembershipListDifferent(receivedMembershipList []*Node) bool {
	var (
		different bool = true
	)

	for remoteNodeIndex := range receivedMembershipList {
		for localNodeIndex := range d.cluster.clusterMetadata {
			switch  {
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
			d.cluser.clusterMetadata = append(d.cluster.clusterMetadata, nodeToJoin)
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
			encodedUpdate, _  := d.marshaler.MarshalSingleNodeUpdate(node.nodeAddress, node.nodeListenPort, node.nodeStatus)
			d.send(fanoutList[index], encodedUpdate)
		}

		// check if the majority quorum is reached.
		if (d.gossipQuorum >= (len(fanoutList) / 2) + 1) {
			break
		}
	}
}

func (d *Dissemination) send(nodeAddress string, gossipMessage []byte) {
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
		if d.ack == 1 {
			d.gossipQuorum += 1
		}
	}
}

func (d *Dissemination) marshalMembershipList(clusterData []*Node) ([]byte, error) {
	var (
		encodedMembershipList bytes.Buffer
		entry []byte
		err error
	)

	for index := range clusterData {
		entry, err = json.Marshal(map[string]any{
			"address": clusterData[index].nodeAddress,
			"port": clusterData[index].listenPort,
			"status": clusterData[index].nodeStatus,
		})

		if err != nil {
			return nil, err
		}

		encodedMembershipList.Write(entry)
	}

	return encodedMembershipList.Bytes(), nil
}