package node

import (
	"context"
	"encoding/json"
	"knucklesdb/store"
	"knucklesdb/swim"
	"knucklesdb/wal"
	"net"
	"strconv"
	"time"

	id "github.com/google/uuid"
)

type Replica struct {
	replicaID        id.UUID
	address          string
	listenPort       string
	kMap             *store.KnucklesMap
	protocolMessages SwimProtocolMessages
	timeoutTime      time.Duration
	swimMarshaler    *swim.ProtocolMarshaer
	clusterJoiner    *swim.ClusterManager
	logger           *wal.ErrorsLogger
	infoLogger       *wal.InfoLogger
	swimGossip       *swim.Dissemination
}

type SwimProtocolMessages struct {
	PiggyBackMsg        swim.PiggyBackMessage
	FailureDetectionMsg swim.DetectionMessage
	JoinRequest         swim.JoinMessage
	SpreadedList        swim.MembershipListMessage
}

type Message struct {
	MethodType string `json:"type"`
	Key        []byte `json:"key"`
	Value      []byte `json:"value,omitempty"`
}

func NewReplica(address string, port string, dbMap *store.KnucklesMap, timeout time.Duration,
	marshaler *swim.ProtocolMarshaer, clusterData *swim.ClusterManager, errLogger *wal.ErrorsLogger,
	infosLog *wal.InfoLogger, dissemination *swim.Dissemination) *Replica {
	return &Replica{
		replicaID:     id.New(),
		address:       address,
		listenPort:    port,
		kMap:          dbMap,
		timeoutTime:   timeout,
		swimMarshaler: marshaler,
		clusterJoiner: clusterData,
		logger:        errLogger,
		infoLogger:    infosLog,
		swimGossip:    dissemination,
	}
}

func (r *Replica) Start() {
	ln, err := net.Listen("tcp", net.JoinHostPort(r.address, r.listenPort))
	if err != nil {
		r.logger.ReportError(err)
	}

	r.infoLogger.ReportInfo("Server Listening")

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
		case "swim", "ping", "piggyback", "membership":
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
		r.kMap.Set(message.Key, message.Value)
		responsePayload, _ = json.Marshal(map[string]string{
			"ack": "1",
		})
	case "get":
		if getErr, value = r.kMap.Get(message.Key); getErr != nil {
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
	case "swim":
		r.HandleSWIMMembershipList(buffer, countBuffer)
	case "ping":
		r.HandlePingSWIMMessage(conn)
	case "piggyback":
		r.HandlePiggyBackSWIMMessage(conn, buffer, countBuffer)
	case "membership":
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
	// prendere il messaggio.
	// ignorarlo se già ricevuto.
	// merge del messaggio.
	// avvio di un nuovo gossip cycle.
}

func (r *Replica) HandleSWIMMembershipList(buffer []byte, bufferLength int) {
	if err := json.Unmarshal(buffer[:bufferLength], r.protocolMessages.SpreadedList); err != nil {
		r.logger.ReportError(err)
		return
	}

	// TODO => Prendere da buffer gli elementi e decodificarli creando una lista di []*swim.Node
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

}
