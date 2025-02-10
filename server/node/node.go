package node

import (
	"encoding/json"
	"fmt"
	"knucklesdb/store"
	"knucklesdb/swim"
	"net"

	id "github.com/google/uuid"
)

type Replica struct {
	replicaID        id.UUID
	address          string
	listenPort       string
	kMap             *store.KnucklesMap
	protocolMessages SwimProtocolMessages
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

func NewReplica(address string, port string, dbMap *store.KnucklesMap) *Replica {
	return &Replica{
		replicaID:  id.New(),
		address:    address,
		listenPort: port,
		kMap:       dbMap,
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
			r.handleSWIMProtocolConnection(conn, buffer, msg.MethodType)
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

func (r *Replica) handleSWIMProtocolConnection(conn net.Conn, buffer []byte, methodType string) {

	switch methodType {
	case "swim":
		r.HandleSWIMFailureDetectionMessage(conn, buffer)
	case "ping":
		r.HandleSWIMPingMessage(conn, buffer)
	case "piggyback":
		r.HandlePiggyBackSWIMMessage(conn, buffer)
	}
}

func (r *Replica) HandleSWIMPingMessage(conn net.Conn, buffer []byte) {

}

func (r *Replica) HandlePiggyBackSWIMMessage(conn net.Conn, buffer []byte) {

}

func (r *Replica) HandleSWIMFailureDetectionMessage(conn net.Conn, buffer []byte) {

}
