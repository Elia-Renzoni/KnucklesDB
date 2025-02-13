package knucklesclient

import (
	"net"
	"time"
	"errors"
	"encoding/json"
)

type KnucklesDBClient struct {
	clientAddr net.IP
	clientListenPort string
	targetNodeAddr net.IP
	targetNodeListenPort string
	replacers []string
	heartbeatingTime time.Duration
	conn net.Conn
}

func NewClient(clientHost, clientListenPort string, targetNodeAddr, targetNodePort string, frequencyTime time.Duration) *KnucklesDBClient {
	return &KnucklesDBClient{
		clientAddr: net.IP(clientHost),
		clientListenPort, clientListenPort,
		targetNodeAddr: net.IP(targetNodeAddr),
		targetNodeListenPort: targetNodePort,
		replacers: make([]string, 0),
		heartbeatingTime: frequencyTime,
	}
}

func (k *KnucklesDBClient) Set(key, value []byte) error {
	if okKey := isEmpty(key); okKey {
		return errors.New("The Key is Empty!")
	}

	if okValue := isEmpty(value); okValue {
		return errors.New("The Value is Empty!")
	}

	joined := net.JoinHostPort(k.targetNodeAddr.String(), k.targetNodeListenPort)
	jsonValue, marshalError := json.Marshal(map[string]string{
		"type": "set",
		"key": key,
		"value": value,
	})

	if marshalError != nil {
		return marshalError
	}

	for {
		time.Sleep(k.heartbeatingTime)
		k.conn, err := net.Dial("tcp", joined)
		if err != nil {
			return err
		}
		k.conn.Write(jsonValue)

		reply := make([]byte, 2024)
		n, _ := k.conn.Read(reply)
		k.conn.Close()
	}
	return nil
}

func (k *KnucklesDBClient) Get(key []byte) ([]byte, error) {
	if okKey := isEmpty(key); okKey {
		return nil, errors.New("The Key is Empty!")
	}

	jsonGetValue, marshalError := json.Marshal(map[string]string{
		"type": "get",
		"key": key,
	})

	if marshalError != nil {
		return nil, marshalError
	}

	joined := net.JoinHostPort(k.targetNodeAddr.String(), k.targetNodeListenPort)

	k.conn, err := net.Dial("tcp", joined)
	if err ! nil {
		return nil, err
	}

	k.conn.Write(jsonGetValue)

	reply := make([]byte, 2024)
	n, _ := k.conn.Read(reply)
	// TODO -> check buffer errors

	k.conn.Close()

	return reply[:n], nil
}

func (k *KnucklesDBClient) SetReplacers(nodes ...string) {

}

func isEmpty(bytesToCheck []byte) bool {
	var resultToReturn bool
	for _, b := range bytesToCheck {
		if b != 0 {
			resultToReturn = true
		} else {
			resultToReturn = false
		}
	}
	return resultToReturn
}
