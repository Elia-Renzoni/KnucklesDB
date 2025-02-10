/**
*	This module contains the implementation of the SWIM protocol according
*	to the paper called 
*   SWIM: Scalable Weakly-consistent Infection-style Process Group Membership Protocol
*   The follow protocol use Round Robin style Algorithm to handle the cluster list.
*   Each swimSchedule time period the goroutine will start a round robin round, during
* 	this round the goroutine will send a ping to the selected node and wait for the ack
*	if the response goes timeout then the parent goroutine will schedule a child one.
*	The child goroutine will handle the piggybacks.
*/
package swim

import (
	"net"
	"time"
	"context"
	"math/rand"
)

type SWIMFailureDetector struct {
	nodesList *ClusterManager
	marshaler *ProtocolMarshaer
	swimMessageAck AckMessage
	kHelperNodes int
	swimSchedule time.Duration
	timeoutTime time.Duration
}

func NewSWIMFailureDetector(nodes *ClusterManager, marshaler *ProtocolMarshaer, helperNodes int, 
	                       sleepTime, timeoutBoundaries time.Duration) *SWIMFailureDetector {
	return &SWIMFailureDetector{
		nodesList: nodes,
		marshaler: marshaler,
		kHelperNodes: helperNodes,
		swimSchedule: sleepTime,
		timeoutTime: timeoutBoundaries,
	}
}

func (s *SWIMFailureDetector) sendPing(nodeHost string, nodeListenPort int) {
	joined := net.JoinHostPort(nodeHost, nodeListenPort)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, s.timeoutTime)
	defer cancel()

	conn, err := net.Dial("tcp", joined)
	defer conn.Close()

	if err != nil {
		// do something... write in WAL
	}

	jsonValue, _ := s.marshaler.MarshalPing()
	conn.Write(jsonValue)

	replyData := make([]byte, 2040)	

	select {
	case <- ctx.Done():
		s.changeNodeState(nodeHost, STATUS_SUSPICIOUS)
		// TODO -> start a gossip cycle
		go s.piggyBack(joined)
	default: 
		count, _ := conn.Read(reply)
		json.Unmarshal(replyData[:count], &s.swimMessageAck)
	}
}

func (s *SWIMFailureDetector) piggyBack(targetInfo string) {	
	if len(s.nodesList.clusterMetadata) < s.helperNodes {
		// TODO -> write error in the WAL
	} else {
		
		var helperResponses []int = make([]int, 0)
		var eliminationCondition bool = true

		for i := 0; i < s.kHelperNodes; i++ {
			randomKHelperNode := rand.Intn(s.kHelperNodes + 1)
			helperNode := s.nodesList.clusterMetadata[randomKHelperNode]
			piggy := s.pingPiggyBack()
			result := piggy(helperNode.nodeAddress.String(), helperNode.nodeListenPort, joined)
			helperResponses = append(helperResponses, result)
		}

		host, _, _ := net.SplitHostPort(targetInfo)

		for _, result := range helperResponses {
			if result == 1 {
				s.changeNodeState(host, STATUS_ALIVE)
				eliminationCondition = false
				break
			}
		}

		if eliminationCondition {
			s.changeNodeState(host, STATUS_REMOVED)
		}
	}
}

func (s *SWIMFailureDetector) pingPiggyBack() func(string, int, string) int {
	return func(parentIP string, parentPort int, targetNode string) int {
		var piggyBackHelperNodeAck AckMessage

		ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
		defer cancel()

		conn, err := net.Dial("tcp", net.JoinHostPort(parentIP, parentPort))
		defer conn.Close()

		if err != nil {
			// TODO -> write error in the WAL
		}

		jsonValue, _ := s.marshaler.MarshalPiggyBack(net.JoinHostPort(parentIP, parentPort), targetNode)
		conn.Write(jsonValue)
		reply := make([]byte, 2040)

		select {
		case <- ctx.Done():
			// 0 mark the target node as not reachable due to
			// the timeout of the helper node.
			return 0
		default: 
			count, _ := conn.Read(reply)
			json.Unmarshal(replyData[:count], &piggyBackHelperNodeAck)
		}
		
		return piggyBackHelperNodeAck.ackContent 
	}
}

func (s *SWIMFailureDetector) replyToPing() {

}

func (s *SWIMFailureDetector) changeNodeState(nodeHost string, nodeUpdatedStatus int) {
	// search and get the node
	for _, node := range s.nodesList.clusterMetadata {
		if node.nodeAddress == nodeHost {
			node.nodeStatus = nodeUpdatedStatus
		}
	}
}

// this method represent the goroutine that has to be called
// by the server
func (s *SWIMFailureDetector) ClusterFailureDetection() {
	for {
		time.Sleep(s.swimSchedule)

		for _, node := range s.nodesList.clusterMetadata {
			s.sendPing(node.nodeAddress.String(), node.nodeListenPort)
		}
	}
}
