package node

import (
	"context"
	"encoding/json"
	"fmt"
	"knucklesdb/store"
	"knucklesdb/swim"
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
}

type SwimProtocolMessages struct {
	PiggyBackMsg        swim.PiggyBackMessage
	FailureDetectionMsg swim.DetectionMessage
	JoinRequest         swim.JoinMessage
}

type Message struct {
	MethodType string `json:"type"`
	Key        []byte `json:"key"`
	Value      []byte `json:"value,omitempty"`
}

func NewReplica(address string, port string, dbMap *store.KnucklesMap, timeout time.Duration,
	marshaler *swim.ProtocolMarshaer, clusterData *swim.ClusterManager) *Replica {
	return &Replica{
		replicaID:     id.New(),
		address:       address,
		listenPort:    port,
		kMap:          dbMap,
		timeoutTime:   timeout,
		swimMarshaler: marshaler,
		clusterJoiner: clusterData,
	}
}

func (r *Replica) Start() {
	ln, err := net.Listen("tcp", net.JoinHostPort(r.address, r.listenPort))
	if err != nil {
		// TODO -> write error message in the WAL
		fmt.Printf("In the replica %s occurred %v", r.replicaID.String(), err)
	}

	fmt.Printf("Server Listening...\n")

	for {
		conn, err := ln.Accept()

		if err != nil {
			// TODO -> write error message in WAL
			fmt.Printf("%v", err)
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

	buffer = make([]byte, 2040)
	n, err := conn.Read(buffer)
	if err != nil {
		// TODO -> write error message in WAL
		fmt.Printf("%v", err)
	}

	if err := json.Unmarshal(buffer[:n], msg); err != nil {
		// TODO -> write to WAL
		fmt.Printf("%v\n", err)
	} else {
		switch msg.MethodType {
		case "swim", "ping", "piggyback":
			r.handleSWIMProtocolConnection(conn, buffer, msg.MethodType, n)
		case "set", "get":
			go r.handleConnection(conn, msg)
		case "join":
			r.handleJoinMembershipMessage(conn, buffer, n)
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
			responsePayload, _ = json.Marshal(map[string][]byte{
				"ack": value,
			})
		}
	}

	conn.Write(responsePayload)
}

func (r *Replica) handleSWIMProtocolConnection(conn net.Conn, buffer []byte, methodType string, countBuffer int) {
	defer conn.Close()

	switch methodType {
	case "swim":
		r.HandleSWIMFailureDetectionMessage(buffer, countBuffer)
	case "ping":
		r.HandlePingSWIMMessage(conn)
	case "piggyback":
		r.HandlePiggyBackSWIMMessage(conn, buffer, countBuffer)
	}
}

func (r *Replica) HandlePingSWIMMessage(conn net.Conn) {
	jsonAckValueToSend, err := r.swimMarshaler.MarshalAckMessage(1)
	if err != nil {
		// Write to WAL.
	}
	conn.Write(jsonAckValueToSend)
}

func (r *Replica) HandlePiggyBackSWIMMessage(conn net.Conn, buffer []byte, bufferLength int) {

	if err := json.Unmarshal(buffer[:bufferLength], &r.protocolMessages.PiggyBackMsg); err != nil {
		// TODO -> Write Error Message to WAL.
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeoutTime)
	defer cancel()

	connHelper, err := net.Dial("tcp", r.protocolMessages.PiggyBackMsg.TargetNode)
	defer conn.Close()

	if err != nil {
		// TODO -> Write to WAL
	}

	jsonEncodedPing, err := r.swimMarshaler.MarshalPing()
	if err != nil {
		// TODO -> Write to WAL
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

// TODO
func (r *Replica) HandleSWIMFailureDetectionMessage(buffer []byte, bufferLength int) {

}

func (r *Replica) handleJoinMembershipMessage(conn net.Conn, buffer []byte, bufferLength int) {
	json.Unmarshal(buffer[:bufferLength], &r.protocolMessages.JoinRequest)
	converted, err := strconv.Atoi(r.protocolMessages.JoinRequest.ListenPort)
	if err != nil {
		bytes, _ := json.Marshal(map[string]any{
			"error": "Malformed Listen Port",
		})
		conn.Write(bytes)
	} else {
		r.clusterJoiner.JoinCluster(net.IP(r.protocolMessages.JoinRequest.IPAddr), converted)
		toSend, _ := r.swimMarshaler.MarshalAckMessage(1)
		conn.Write(toSend)
	}

	conn.Close()
}