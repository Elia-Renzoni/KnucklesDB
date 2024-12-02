package node

import (
	"net"
	id "github.com/google/uuid"
	"context"
	"encoding/json"
)

type Replica struct {
	replicaID id.UUID 
	address string
	listenPort int
	quitError chan <-struct{}
}

type Message struct {
	id int `json:"id"`
	methodType string `json:"type"`	
	methodName string `json:"method"`   
	parameter string `json:"parameter"`
}

func NewReplica(address string, port int, errorsChannel chan <-struct{}) *Replica {
	return &Replica{
		replicaID: id.New(),
		address: address,
		listenPort: port,
		quitError: errorsChannel,
	}
}

func (r *Replica) Start() {
	ln, err := net.Listen("tcp", r.listenPort)
	if err != nil {
		fmt.Printf("In the replica %s occurred %v", r.replicaID.String(), err)
	}	

	for {
		conn, err := ln.Accept()
		if err != nil {
			// to handle
		}

		go handleConnection(conn)
	}
}

func handleConnection(ctx context.Context, conn net.Conn) error {
	defer conn.Close()

	messageBuffer := make([]byte, 2024)
	_, err := conn.Read(messageBuffer)
	if err != nil {
		return err
	}

	var msg = &Message{}

	json.Unmarshal(msg, messageBuffer)
	switch msg.methodType {
	case "set":
		handleSetRequest()
	case "get";
		handleGetRequest()
	default:
		return errors.New("Invalid Method Type")
	}
}

func handleSetRequest() {

}

func handleGetRequest() {

}