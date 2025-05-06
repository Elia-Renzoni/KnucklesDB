package node

import (
	"context"
	"encoding/json"
	"knucklesdb/consensus"
	"knucklesdb/store"
	"knucklesdb/swim"
	"knucklesdb/vvector"
	"knucklesdb/wal"
	"net"
	"strconv"
	"time"
	"sync"
	"fmt"

	id "github.com/google/uuid"
)

type Replica struct {
	replicaID            id.UUID
	host, port string
	kMap                 *store.KnucklesMap
	protocolMessages     SwimProtocolMessages
	versionVectorMessage consensus.PipelinedMessage
	timeoutTime          time.Duration
	swimMarshaler        *swim.ProtocolMarshaer
	clusterJoiner        *swim.ClusterManager
	logger               *wal.ErrorsLogger
	infoLogger           *wal.InfoLogger
	swimGossip           *swim.Dissemination
	gossipConsensus      *consensus.Gossip
	versionVectorUtils   *vvector.DataVersioning
	syncronizer *sync.WaitGroup
}

type SwimProtocolMessages struct {
	PiggyBackMsg        swim.PiggyBackMessage
	FailureDetectionMsg swim.DetectionMessage
	JoinRequest         swim.JoinMessage
	SpreadedList        swim.MembershipListMessage
	NodeUpdate          swim.SWIMUpdateMessage
}

type Message struct {
	MethodType string `json:"type"`
	Key        []byte `json:"key"`
	Value      []byte `json:"value,omitempty"`
}

func NewReplica(host, port string, uuid id.UUID, dbMap *store.KnucklesMap, timeout time.Duration,
	marshaler *swim.ProtocolMarshaer, clusterData *swim.ClusterManager, errLogger *wal.ErrorsLogger,
	infosLog *wal.InfoLogger, dissemination *swim.Dissemination, gossip *consensus.Gossip,
	versionVector *vvector.DataVersioning, syncronizer *sync.WaitGroup) *Replica {
	return &Replica{
		replicaID:          uuid,
		host: host,
		port: port,
		kMap:               dbMap,
		timeoutTime:        timeout,
		swimMarshaler:      marshaler,
		clusterJoiner:      clusterData,
		logger:             errLogger,
		infoLogger:         infosLog,
		swimGossip:         dissemination,
		gossipConsensus:    gossip,
		versionVectorUtils: versionVector,
		syncronizer: syncronizer,
	}
}

func (r *Replica) Start() {
	ln, err := net.Listen("tcp", net.JoinHostPort(r.host, r.port))
	if err != nil {
		r.logger.ReportError(err)
	}

	r.infoLogger.ReportInfo("Server Listening")

	// make possibile the Join 
	r.syncronizer.Done()
	r.infoLogger.ReportInfo("Done!")
	for {
		conn, err := ln.Accept()

		if err != nil {
			r.logger.ReportError(err)
		}

		r.serveRequest(conn)
	}
}

func (r *Replica) serveRequest(conn net.Conn) {
	var (
		buffer []byte
		msg    = &Message{}
		ctx    = context.Background()
	)

	_, cancel := context.WithDeadline(ctx, time.Time{}.Add(10*time.Second))
	defer cancel()

	buffer = make([]byte, 5040)
	n, err := conn.Read(buffer)
	if err != nil {
		r.logger.ReportError(err)
	}


	if err := json.Unmarshal(buffer[:n], msg); err != nil {
		r.logger.ReportError(err)
	} else {
		switch msg.MethodType {
		case "swim-update", "ping", "piggyback", "membership":
			r.handleSWIMProtocolConnection(conn, buffer, msg.MethodType, n)
		case "set", "get":
			go r.handleConnection(conn, msg)
		case "join":
			r.handleJoinMembershipMessage(conn, buffer, n)
		case "gossip":
			r.handleConsensusAgreementMessage(conn, buffer, n)
		default:
			toSend, _ := json.Marshal(map[string]string{
				"error": "Illegal Method Type",
			})
			conn.Write(toSend)
		}
	}
}

func (r *Replica) handleConnection(conn net.Conn, message *Message) {
	var (
		getErr          error
		value           []byte
		responsePayload []byte
	)

	defer conn.Close()

	switch message.MethodType {
	case "set":
		r.kMap.Set(message.Key, message.Value, 0)
		responsePayload, _ = json.Marshal(map[string]string{
			"ack": "1",
		})
	case "get":
		if getErr, value, _ = r.kMap.Get(message.Key); getErr != nil {
			responsePayload, _ = json.Marshal(map[string]string{
				"ack": getErr.Error(),
			})
		} else {
			responsePayload, _ = json.Marshal(map[string]string{
				"ack": string(value),
			})
		}
	}

	conn.Write(responsePayload)
}

func (r *Replica) handleSWIMProtocolConnection(conn net.Conn, buffer []byte, methodType string, countBuffer int) {
	defer conn.Close()

	switch methodType {
	case "membership":
		r.HandleSWIMMembershipList(conn, buffer, countBuffer)
	case "ping":
		r.HandlePingSWIMMessage(conn)
	case "piggyback":
		r.HandlePiggyBackSWIMMessage(conn, buffer, countBuffer)
	case "swim-update":
		r.HandleSWIMGossipMessage(conn, buffer, countBuffer)
	}
}

func (r *Replica) HandlePingSWIMMessage(conn net.Conn) {
	r.infoLogger.ReportInfo("Ping Message Arrived")
	jsonAckValueToSend, err := r.swimMarshaler.MarshalAckMessage(1)
	if err != nil {
		r.logger.ReportError(err)
	}
	conn.Write(jsonAckValueToSend)
}

func (r *Replica) HandlePiggyBackSWIMMessage(conn net.Conn, buffer []byte, bufferLength int) {
	r.infoLogger.ReportInfo("PiggyBack Message Arrived")

	if err := json.Unmarshal(buffer[:bufferLength], &r.protocolMessages.PiggyBackMsg); err != nil {
		r.logger.ReportError(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeoutTime)
	defer cancel()

	connHelper, err := net.Dial("tcp", r.protocolMessages.PiggyBackMsg.TargetNode)

	if err != nil {
		r.logger.ReportError(err)
		jsonValueNeg, _ := r.swimMarshaler.MarshalAckMessage(0)
		conn.Write(jsonValueNeg)
		return
	}
	defer connHelper.Close()

	jsonEncodedPing, err := r.swimMarshaler.MarshalPing()
	if err != nil {
		r.logger.ReportError(err)
	}

	connHelper.Write(jsonEncodedPing)

	select {
	// f the piggyback transmission to the target node times out,
	// I must send a message to the parent node containing a negative acknowledgment (ACK) with a value of 0.
	case <-ctx.Done():
		jsonValueNeg, _ := r.swimMarshaler.MarshalAckMessage(0)
		conn.Write(jsonValueNeg)
	default:
		jsonValuePos, _ := r.swimMarshaler.MarshalAckMessage(1)
		conn.Write(jsonValuePos)
	}
}

func (r *Replica) HandleSWIMGossipMessage(conn net.Conn, buffer []byte, bufferLength int) {
	if err := json.Unmarshal(buffer[:bufferLength], &r.protocolMessages.NodeUpdate); err != nil {
		r.logger.ReportError(err)
		return
	}

	r.infoLogger.ReportInfo("SWIM Node Update Arrived via Gossip")

	node := swim.NewNode(r.protocolMessages.NodeUpdate.NodeAddress, r.protocolMessages.NodeUpdate.NodeListenPort, r.protocolMessages.NodeUpdate.NodeStatus)
	if different := r.swimGossip.IsUpdateDifferent(node); different {
		r.swimGossip.MergeUpdates(node)
		fanoutList := r.clusterJoiner.SetFanoutList()
		r.infoLogger.ReportInfo("SWIM FANOUT FOR SPREADING UPDATES")
		go r.swimGossip.SpreadMembershipListUpdates(fanoutList, node)

		jsonAck, _ := r.swimMarshaler.MarshalAckMessage(1)
		conn.Write(jsonAck)
	} else {
		jsonAck, _ := r.swimMarshaler.MarshalAckMessage(0)
		conn.Write(jsonAck)
	}
}

func (r *Replica) HandleSWIMMembershipList(conn net.Conn, buffer []byte, bufferLength int) {
	var (
		seed bool = true
		fanoutList []string
	)

	if err := json.Unmarshal(buffer[:bufferLength], &r.protocolMessages.SpreadedList); err != nil {
		r.logger.ReportError(err)
		return
	}

	r.infoLogger.ReportInfo("Broadcast MembershipList Arrived")

	decodedMembershipList := r.swimGossip.TransformMembershipList(r.protocolMessages.SpreadedList)
	if isDifferent := r.swimGossip.IsMembershipListDifferent(decodedMembershipList); isDifferent {
		r.swimGossip.MergeMembershipList(decodedMembershipList)

		fmt.Printf("Sender Address -----> %s \n", r.protocolMessages.SpreadedList.SenderAddr)
		if result := r.clusterJoiner.CheckIfFanoutIsPossible(r.protocolMessages.SpreadedList.SenderAddr, net.JoinHostPort(r.host, r.port)); result {
			r.infoLogger.ReportInfo("FANOUT possiible")
			fanoutList = r.clusterJoiner.SetFanoutList()
			for r.checkFanoutList(r.protocolMessages.SpreadedList.SenderAddr, fanoutList) {
				fanoutList = r.clusterJoiner.SetFanoutList()
			}

			fmt.Println("Fanout List: ")
			for _ , node := range fanoutList {
				fmt.Printf("%s", node)
			}
		
			if r.port != "5050" {
				seed = false
			}

			go r.swimGossip.SpreadMembershipList(decodedMembershipList, fanoutList, seed)
		} else {
			r.infoLogger.ReportInfo("Non è possibile fare il FANOUT")
		}

		jsonAck, _ := r.swimMarshaler.MarshalAckMessage(1)
		conn.Write(jsonAck)
	} else {
		jsonAck, _ := r.swimMarshaler.MarshalAckMessage(0)
		conn.Write(jsonAck)
	}
}

func (r *Replica) checkFanoutList(remoteAddr string, calculatedAddrs []string) bool {
	var result bool

	for _, fanoutNode := range calculatedAddrs {
		switch {
		case remoteAddr == fanoutNode, fanoutNode == net.JoinHostPort(r.host, r.port), fanoutNode == "127.0.0.1:5050":
			result = true
		}
	}

	return result
}

func (r *Replica) handleJoinMembershipMessage(conn net.Conn, buffer []byte, bufferLength int) {
	r.infoLogger.ReportInfo("Join Message Arrived")

	json.Unmarshal(buffer[:bufferLength], &r.protocolMessages.JoinRequest)
	converted, err := strconv.Atoi(r.protocolMessages.JoinRequest.ListenPort)
	if err != nil {
		bytes, _ := json.Marshal(map[string]any{
			"error": "Malformed Listen Port",
		})
		conn.Write(bytes)
	} else {
		r.clusterJoiner.JoinCluster(r.protocolMessages.JoinRequest.IPAddr, converted)
		toSend, _ := r.swimMarshaler.MarshalAckMessage(1)
		conn.Write(toSend)
	}

	conn.Close()
}

func (r *Replica) handleConsensusAgreementMessage(conn net.Conn, messageBuffer []byte, messageBufferLength int) {
	defer conn.Close()
	// unmarshal the message received by peers via gossip.
	json.Unmarshal(messageBuffer[:messageBufferLength], &r.versionVectorMessage)
	ackMessage, _ := r.swimMarshaler.MarshalAckMessage(1)

	if ok, clock := r.gossipConsensus.SearchReplica(r.versionVectorMessage.ReplicaUUID); !ok {
		// if the replica is not present in the termination hash map is considered a new message from the
		// replica.
		r.gossipConsensus.AddReplicaInTerminationMap(r.versionVectorMessage.ReplicaUUID, r.versionVectorMessage.LogicalClock)
	} else {
		// new message from nodes
		if clock != r.versionVectorMessage.LogicalClock {
			// process the message
			// and then forward the message
			r.gossipConsensus.AddReplicaInTerminationMap(r.versionVectorMessage.ReplicaUUID, r.versionVectorMessage.LogicalClock)
			// performing a LLW between the received pipeline of messages
			r.gossipConsensus.PipelinedLLW(r.versionVectorMessage.Pipeline)

			// perfoming a LLW between the received pipeline and the memory content.
			r.performLLW(r.versionVectorMessage.Pipeline)

			// start a new gossip round
			go func() {
				fanoutList := r.clusterJoiner.SetFanoutList()
				for nodeIndex := range fanoutList {
					r.gossipConsensus.Send(fanoutList[nodeIndex], messageBuffer[:messageBufferLength])
				}
			}()
		}
		conn.Write(ackMessage)
	}
}

func (r *Replica) performLLW(pipeline []vvector.VersionVectorMessage) {
	for pipelineNodeIndex := range pipeline {
		err, _, inMemoryVersion := r.kMap.Get(pipeline[pipelineNodeIndex].Key)
		if err != nil {
			// the value is not in memory, we need to perform the first set
			r.kMap.Set(pipeline[pipelineNodeIndex].Key, pipeline[pipelineNodeIndex].Value, pipeline[pipelineNodeIndex].Version)
		} else {
			// if the value is already in the buffer pool we need to confront the
			// versions to get a correct LLW.
			r.versionVectorUtils.CompareAndUpdateVersions(pipeline[pipelineNodeIndex], inMemoryVersion)
			switch r.versionVectorUtils.Order {
			case vvector.HAPPENS_AFTER:
				// if the received version is greater then update the memorized version
				// otherwise the memorized version is more updated.
				r.kMap.Set(pipeline[pipelineNodeIndex].Key, pipeline[pipelineNodeIndex].Value, pipeline[pipelineNodeIndex].Version)
			}
		}
	}
}
