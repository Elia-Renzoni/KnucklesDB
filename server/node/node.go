package node

import (
	"net"
	id "github.com/google/uuid"
	"encoding/json"
	"knucklesdb/clock"
	"knucklesdb/store"
	"strings"
)

type Replica struct {
	replicaID id.UUID 
	address string
	listenPort string
	internalClock *clock.KnucklesClock
	db *store.KnucklesDB
	values *store.DBvalues
}

type Message struct {
	methodType string `json:"type"`	
	methodName string `json:"method"`   
	parameter string `json:"parameter"`
}

func NewReplica(address string, port string, logiclaClock *clock.KnucklesClock,
				db *store.KnucklesDB, values *store.DBvalues) *Replica {
	return &Replica{
		replicaID: id.New(),
		address: address,
		listenPort: port,
		internalClock: logicalClock,
		db: db,
		values: values,
	}
}

func (r *Replica) Start() {
	ln, err := net.Listen("tcp", net.JoinHostPort(r.address, r.listenPort))
	if err != nil {
		fmt.Printf("In the replica %s occurred %v", r.replicaID.String(), err)
	}	

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
		return err
	}

	var msg = &Message{}

	json.Unmarshal(msg, messageBuffer)
	switch msg.methodType {
	case "set":
		if setErr = handleSetRequest(msg.methodName, msg.parameter); setErr != nil {
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
	case "get";
		if getErr, value = handleGetRequest(msg.methodName, msg.parameter); getErr != nil {
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
			"error": "Illegal Method Type"
		})
	}
}

func (r *Replica) handleSetRequest(methodName string, parameter string) error {
	switch methodName {
	case "ip":
		host, port, _ := net.SplitHostPort(parameter)
		value := store.NewDBValues(host, port, r.internalClock.GetLogicalClock(), "")
		r.db.SetWithIpAddressOnly(host, value)
	case "end":
		splitted := strings.Split(parameter, ":")
		value := store.NewDBValues(nil, splitted[1], r.internalClock.GetLogicalClock(), splitted[0])
		r.db.SetWithEndpointOnly(splitted[0], value)
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
		return errors.New("Illegal Parameter")
	}
	return nil, value
}