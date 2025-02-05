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
)

type SWIMFailureDetector struct {
	nodesList *ClusterManager
	kHelperNodes int
	swimSchedule time.Time
	timeoutTime time.Time
}

func NewSWIMFailureDetector(nodes *ClusterManager, helperNodes int, sleepTime, timeoutBoundaries time.Time) *SWIMFailureDetector {
	return &SWIMFailureDetector{
		nodesList: nodes,
		kHelperNodes: helperNodes,
		swimSchedule: sleepTime,
		timeoutTime: timeoutBoundaries,
	}
}

func (s *SWIMFailureDetector) sendPing() {

}

func (s *SWIMFailureDetector) piggyBack() {

}

// this method represent the goroutine that has to be called
// by the server
func (s *SWIMFailureDetector) ClusterFailureDetection() {

}
