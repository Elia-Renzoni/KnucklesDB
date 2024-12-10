package node

import (
	"net"
	id "github.com/google/uuid"
	"encoding/json"
	"knucklesdb/clock"
	"knucklesdb/store"
	"fmt"
	"errors"
)

type Replica struct {
	replicaID id.UUID 
	address string
	listenPort string
	internalClock *clock.LogicalClock
	db *store.KnucklesDB
}

type Message struct {
	MethodType string `json:"type"`	
	MethodName string `json:"method"`   
	Parameter string `json:"parameter"`
	Port int `json:"port"`
}

func NewReplica(address string, port string, logicalClock *clock.LogicalClock,
				db *store.KnucklesDB) *Replica {
	return &Replica{
		replicaID: id.New(),
		address: address,
		listenPort: port,
		internalClock: logicalClock,
		db: db,
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
		setErr error
		getErr error
		value *store.DBvalues
		toWrite string
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

	switch msg.MethodName {
	case "set":
		if setErr = r.handleSetRequest(msg.MethodType, msg.Parameter, msg.Port); setErr != nil {
			responsePayload, _ = json.Marshal(map[string]string{
				"error": setErr.Error(),
			})
		} else {
			responsePayload, _ = json.Marshal(map[string]string{
				"ack": "1",
			})		
		}
	case "get": 
		if getErr, value = r.handleGetRequest(msg.MethodType, msg.Parameter); getErr != nil {
			responsePayload, _ = json.Marshal(map[string]string{
				"error": getErr.Error(),
			})
		} else {
			if value.GetIpAddress() != nil {
				toWrite = value.GetIpAddress().String()
			} else {
				toWrite = value.GetOptionalEndpoint()
			}

			responsePayload, _ = json.Marshal(map[string]string{
				"ack": toWrite,
			})
		}
	default:
		responsePayload, _ = json.Marshal(map[string]string{
			"error": "Illegal Method Type",
		})
	}

	conn.Write(responsePayload)
}

func (r *Replica) handleSetRequest(methodType string, parameter string, port int) error {
	switch methodType {
	case "ip":
		value := store.NewDBValues(net.ParseIP(parameter), port, r.internalClock.GetLogicalClock(), "")
		r.db.SetWithIpAddressOnly(parameter, value)
	case "end":
		value := store.NewDBValues(nil, port, r.internalClock.GetLogicalClock(), parameter)
		r.db.SetWithEndpointOnly(parameter, value)
	default:
		return errors.New("Illegal Parameterr")
	}
	return nil
}

func (r *Replica) handleGetRequest(methodType string, parameter string) (error, *store.DBvalues) {
	var (
		err error
		value *store.DBvalues
	)

	switch methodType {
	case "ip":
		value, err = r.db.SearchWithIpOnly(parameter)
		if err != nil && value == nil {
			return err, nil
		}
	case "end":
		value, err = r.db.SearchWithEndpointOnly(parameter)
		if err != nil && value == nil {
			return err, nil
		}
	default:
		return errors.New("Illegal Parameter"), nil
	}
	return nil, value
}