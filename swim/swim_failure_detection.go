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

	conn, err := net.DialContext(ctx, "tcp", joined)
	defer conn.Close()

	if err != nil {
		// do something...
	}

	jsonValue, _ := s.marshaler.MarshalPing()
	conn.Write(jsonValue)

	replyData := make([]byte, 2040)
	count, _ := conn.Read(reply)
	json.Unmarshal(replyData[:count], &s.swimMessageAck)

	select {
	case <- ctx.Done():
		// do something...
	}
}

func (s *SWIMFailureDetector) piggyBack() {

}

func (s *SWIMFailureDetector) changeNodeState() {

}

// this method represent the goroutine that has to be called
// by the server
func (s *SWIMFailureDetector) ClusterFailureDetection() {

}
