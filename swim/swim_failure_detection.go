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
	"context"
	"encoding/json"
	"math/rand"
	"net"
	"time"
	"fmt"
	"strconv"
	"syscall"
	"os"
	"knucklesdb/wal"
)

type SWIMFailureDetector struct {

	// used for the cluster metadata list
	nodesList *ClusterManager

	// message marshaler
	marshaler      *ProtocolMarshaer
	swimMessageAck AckMessage

	// number of the K helper nodes to use
	// during the piggyback session
	kHelperNodes int

	// time for scheduling the swim protocol session
	swimSchedule time.Duration

	// timeout time
	timeoutTime time.Duration

	logger *wal.InfoLogger
}

func NewSWIMFailureDetector(nodes *ClusterManager, marshaler *ProtocolMarshaer, helperNodes int,
	sleepTime, timeoutBoundaries time.Duration, logger *wal.InfoLogger) *SWIMFailureDetector {
	return &SWIMFailureDetector{
		nodesList:    nodes,
		marshaler:    marshaler,
		kHelperNodes: helperNodes,
		swimSchedule: sleepTime,
		timeoutTime:  timeoutBoundaries,
		logger: logger,
	}
}

/*
*	@brief this method send a ping to the target node
*	@param IP address of the target node
*	@param listen port of the target node
 */
func (s *SWIMFailureDetector) sendPing(nodeHost string, nodeListenPort int) {
	joined := net.JoinHostPort(nodeHost, strconv.Itoa(nodeListenPort))
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, s.timeoutTime)
	defer cancel()

	conn, err := net.Dial("tcp", joined)

	if err != nil {
		// do something... write in WAL
		if opErr, ok := err.(*net.OpError); ok {
			if sysErr, okErr := opErr.Err.(*os.SyscallError); okErr {
				if sysErr.Err == syscall.ECONNREFUSED {
					s.changeNodeState(nodeHost, strconv.Itoa(nodeListenPort), STATUS_SUSPICIOUS)
					go s.piggyBack(joined)
					return
				}
			}
		}
	}

	defer conn.Close()

	jsonValue, _ := s.marshaler.MarshalPing()
	conn.Write(jsonValue)

	replyData := make([]byte, 2040)

	select {
	// timeout occured
	case <-ctx.Done():
		s.changeNodeState(nodeHost, strconv.Itoa(nodeListenPort), STATUS_SUSPICIOUS)
		// TODO -> start a gossip cycle
		go s.piggyBack(joined)
	default:
		count, _ := conn.Read(replyData)
		json.Unmarshal(replyData[:count], &s.swimMessageAck)
		fmt.Println(s.swimMessageAck)
	}
}

/*
*	@brief this method implements the piggy back logics of the swim protocol
*	the parent node sends to the K helper nodes (chosed randomly) a message
*	indicating ther target node to ping.
*	@param target address and listen port
 */
func (s *SWIMFailureDetector) piggyBack(targetInfo string) {
	if len(s.nodesList.clusterMetadata) < s.kHelperNodes {
		// TODO -> write error in the WAL
	} else {

		var helperResponses []int = make([]int, 0)
		var eliminationCondition bool = true

		for i := 0; i < s.kHelperNodes; i++ {
			randomKHelperNode := rand.Intn(s.kHelperNodes)

			// select the K node
			helperNode := s.nodesList.clusterMetadata[randomKHelperNode]

			piggy := s.pingPiggyBack()

			// send the piggyback message indicating the target node address and the
			// parent address
			result := piggy(helperNode.nodeAddress, helperNode.nodeListenPort, targetInfo)

			// store the result of the piggyback operation
			helperResponses = append(helperResponses, result)
		}

		host, port, _ := net.SplitHostPort(targetInfo)

		for _, result := range helperResponses {
			// if there is just only a 1 in the results
			// the target node is considered alive
			if result == 1 {
				s.changeNodeState(host, port, STATUS_ALIVE)
				eliminationCondition = false
				break
			}
		}

		// if every result are made of 0 then
		// the target node must be considered removed
		if eliminationCondition {
			s.changeNodeState(host, port, STATUS_REMOVED)
			s.nodesList.DeleteNodeFromCluster(host, port)
		}

		// only for testing
		for _, v := range s.nodesList.clusterMetadata {
			fmt.Printf("Node ID = %s %d Node Status: %d \n", v.nodeAddress, v.nodeListenPort, v.nodeStatus)
		}
	}
}

func (s *SWIMFailureDetector) pingPiggyBack() func(string, int, string) int {
	return func(parentIP string, parentPort int, targetNode string) int {
		var piggyBackHelperNodeAck AckMessage

		ctx, cancel := context.WithTimeout(context.Background(), s.timeoutTime)
		defer cancel()

		castedPort := strconv.Itoa(parentPort)
		conn, err := net.Dial("tcp", net.JoinHostPort(parentIP, castedPort))

		if err != nil {
			// TODO -> write error in the WAL.
			return 0
		}
		defer conn.Close()

		jsonValue, _ := s.marshaler.MarshalPiggyBack(net.JoinHostPort(parentIP, string(parentPort)), targetNode)
		conn.Write(jsonValue)
		reply := make([]byte, 2040)

		select {
		case <-ctx.Done():
			// 0 mark the target node as not reachable due to
			// the timeout of the helper node.
			return 0
		default:
			count, _ := conn.Read(reply)
			json.Unmarshal(reply[:count], &piggyBackHelperNodeAck)
		}

		return piggyBackHelperNodeAck.AckContent
	}
}

func (s *SWIMFailureDetector) changeNodeState(nodeHost, nodePort string, nodeUpdatedStatus int) {
	castedPort, _ := strconv.Atoi(nodePort)

	// search and get the node
	for _, node := range s.nodesList.clusterMetadata {
		if node.nodeAddress == nodeHost && node.nodeListenPort == castedPort {
			node.nodeStatus = nodeUpdatedStatus
		}
	}
}

// this method represent the goroutine that has to be called
// by the server
func (s *SWIMFailureDetector) ClusterFailureDetection() {
	s.logger.ReportInfo("SWIM Protocol On")
	for {
		time.Sleep(s.swimSchedule)

		for _, node := range s.nodesList.clusterMetadata {
			if node != nil {
				if node.nodeStatus != STATUS_REMOVED {
					s.sendPing(node.nodeAddress, node.nodeListenPort)
				}
			}
			
		}
	}
}
