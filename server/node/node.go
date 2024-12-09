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
	//values *store.DBvalues
}

type Message struct {
	methodType string `json:"type"`	
	methodName string `json:"method"`   
	parameter string `json:"parameter"`
	port int `json:"port"`
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
	)

	defer conn.Close()

	messageBuffer := make([]byte, 2024)
	_, err := conn.Read(messageBuffer)
	if err != nil {
		fmt.Printf("%v", err)
	}

	var msg = &Message{}

	json.Unmarshal(messageBuffer, msg)
	switch msg.methodType {
	case "set":
		if setErr = r.handleSetRequest(msg.methodName, msg.parameter, msg.port); setErr != nil {
			payload, _ := json.Marshal(map[string]string{
				"error": err.Error(),
			})
			conn.Write(payload)
		} else {
			payload, _ := json.Marshal(map[string]string{
				"ack": "1",
			})		
			conn.Write(payload)
		}
	case "get": 
		if getErr, value = r.handleGetRequest(msg.methodName, msg.parameter); getErr != nil {
			payload, _ := json.Marshal(map[string]string{
				"error": err.Error(),
			})
			conn.Write(payload)
		} else {
			if value.GetIpAddress() != nil {
				toWrite = value.GetIpAddress().String()
			} else {
				toWrite = value.GetOptionalEndpoint()
			}

			payload, _ := json.Marshal(map[string]string{
				"ack": toWrite,
			})
			conn.Write(payload)
		}
	default:
		payload, _ := json.Marshal(map[string]string{
			"error": "Illegal Method Type",
		})
		conn.Write(payload)
	}
}

func (r *Replica) handleSetRequest(methodName string, parameter string, port int) error {
	switch methodName {
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

func (r *Replica) handleGetRequest(methodName string, parameter string) (error, *store.DBvalues) {
	var (
		err error
		value *store.DBvalues
	)

	switch methodName {
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