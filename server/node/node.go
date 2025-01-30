package node

import (
	"net"
	id "github.com/google/uuid"
	"encoding/json"
	"knucklesdb/store"
	"fmt"
)

type Replica struct {
	replicaID id.UUID 
	address string
	listenPort string
	kMap *store.KnucklesMap
}

type Message struct {
	MethodType string `json:"type"`	
	Key []byte `json:"key"`
	Value []byte `json:"value,omitempty"`
}

func NewReplica(address string, port string, dbMap *store.KnucklesMap) *Replica {
	return &Replica{
		replicaID: id.New(),
		address: address,
		listenPort: port,
		kMap: dbMap,
	}
}

func (r *Replica) Start() {
	ln, err := net.Listen("tcp", net.JoinHostPort(r.address, r.listenPort))
	if err != nil {
		fmt.Printf("In the replica %s occurred %v", r.replicaID.String(), err)
	}	

	fmt.Printf("Server Listening...\n")

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Printf("%v", err)
		}

		go r.handleConnection(conn)
	}
}

func (r *Replica) handleConnection(conn net.Conn) {
	var (
		getErr error
		//toWrite string
		value []byte
		responsePayload []byte
		msg = &Message{}
	)

	defer conn.Close()

	messageBuffer := make([]byte, 2024)
	n, err := conn.Read(messageBuffer)
	if err != nil {
		fmt.Printf("%v", err)
	}

	fmt.Printf(string(messageBuffer[:n]))

	if err := json.Unmarshal(messageBuffer[:n], msg); err != nil {
		fmt.Printf("\n%v\n", err)
	}

	switch msg.MethodType {
	case "set":
		r.kMap.Set(msg.Key, msg.Value)
		responsePayload, _ = json.Marshal(map[string]string{
			"ack": "1",
		})
	case "get": 
		if getErr, value = r.kMap.Get(msg.Key); getErr != nil {
			responsePayload, _ = json.Marshal(map[string]string{
				"error": getErr.Error(),
			})
		} else {
			//toWrite = string(value)

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