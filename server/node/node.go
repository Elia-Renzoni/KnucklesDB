package node

import (
	"encoding/json"
	"fmt"
	"knucklesdb/store"
	"knucklesdb/swim"
	"net"
	"context"
	"time"
	id "github.com/google/uuid"
)

type Replica struct {
	replicaID        id.UUID
	address          string
	listenPort       string
	kMap             *store.KnucklesMap
	protocolMessages SwimProtocolMessages
	timeoutTime time.Duration
	swimMarshaler *swim.ProtocolMarshaer
}

type SwimProtocolMessages struct {
	PiggyBackMsg        swim.PiggyBackMessage
	FailureDetectionMsg swim.DetectionMessage
}

type Message struct {
	MethodType string `json:"type"`
	Key        []byte `json:"key"`
	Value      []byte `json:"value,omitempty"`
}

func NewReplica(address string, port string, dbMap *store.KnucklesMap, timeout time.Duration,
	           marshaler *swim.ProtocolMarshaer) *Replica {
	return &Replica{
		replicaID:  id.New(),
		address:    address,
		listenPort: port,
		kMap:       dbMap,
		timeoutTime: timeout,
		swimMarshaler: marshaler,
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
	)

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
		default:
			go r.handleConnection(conn, msg)
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
				"error": getErr.Error(),
			})
		} else {
			responsePayload, _ = json.Marshal(map[string][]byte{
				"ack": value,
			})
		}
	default:
		responsePayload, _ = json.Marshal(map[string]string{
			"error": "Illegal Method Type",
		})
	}

	conn.Write(responsePayload)
}

func (r *Replica) handleSWIMProtocolConnection(conn net.Conn, buffer []byte, methodType string, countBuffer int) {
	switch methodType {
	case "swim":
		r.HandleSWIMFailureDetectionMessage(conn, buffer, countBuffer)
	case "ping":
		r.HandleSWIMPingMessage(conn, buffer, countBuffer)
	case "piggyback":
		r.HandlePiggyBackSWIMMessage(conn, buffer, countBuffer)
	}
}

func (r *Replica) HandleSWIMPingMessage(conn net.Conn, buffer []byte, bufferLength int) {

}

func (r *Replica) HandlePiggyBackSWIMMessage(conn net.Conn, buffer []byte, bufferLength int) {
	var piggyBackRequest r.protocolMessages.PiggyBackMsg

	if err := json.Unmarshal(buffer[:bufferLength], &piggyBackRequest); err != nil {
		// TODO -> Write Error Message to WAL.
	}

	host, port, _ := net.SplitHostPort(piggyBackRequest.TargetNode)
	ctx, cancel := context.WithTimeout(context.Background(), r.timeoutTime)
	defer cancel() 

	conn, err := net.Dial("tcp", piggyBackRequest.TargetNode)
	defer conn.Close()

	if err != nil {
		// TODO -> Write to WAL
	}

	jsonEncodedPing, err := r.swimMarshaler.MarshalPing()
	if err != nil {
		// TODO -> Write to WAL
	}

	conn.Write(jsonEncodedPing)

	select {
	// f the piggyback transmission to the target node times out, 
	// I must send a message to the parent node containing a negative acknowledgment (ACK) with a value of 0.
	case <- ctx.Done():
		r.writeBackToParentNode(0, piggyBackRequest.ParentNode)
	default:
		r.writeBackToParentNode(1, piggyBackRequest.ParentNode)
	}
}

func (r *Replica) HandleSWIMFailureDetectionMessage(conn net.Conn, buffer []byte, bufferLength int) {

}

func (r *Replica) writeBackToParentNode(pingAckValue int, parentNodeInfos string) {
	conn, err := net.Dial("tcp", parentNodeInfos)
	defer conn.Close()

	jsonValueToSend, _ := r.swimMarshaler.MarshalAckMessage(pingAckValue)
	conn.Write(jsonValueToSend)
}
