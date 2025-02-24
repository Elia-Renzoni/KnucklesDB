package client

import (
	"net"
	"time"
	"errors"
	"encoding/json"
	"fmt"
)

type KnucklesDBClient struct {
	clientAddr string
	clientListenPort string
	targetNodeAddr string
	targetNodeListenPort string
	replacers []string
	heartbeatingTime time.Duration
	conn net.Conn
}

type ServerMessages struct {
	Ack string `json:"ack"`
}


func NewClient(clientHost, clientListenPort string, targetNodeAddr, targetNodePort string, frequencyTime time.Duration) *KnucklesDBClient {
	return &KnucklesDBClient{
		clientAddr: clientHost,
		clientListenPort: clientListenPort,
		targetNodeAddr: targetNodeAddr,
		targetNodeListenPort: targetNodePort,
		replacers: make([]string, 0),
		heartbeatingTime: frequencyTime,
	}
}

func (k *KnucklesDBClient) Set(key, value []byte) error {
	var response ServerMessages

	fmt.Printf("yooo")

	/*if okKey := isEmpty(key); okKey {
		return errors.New("The Key is Empty!")
	}

	if okValue := isEmpty(value); okValue {
		return errors.New("The Value is Empty!")
	}*/

	fmt.Printf("yuuuuu")

	joined := net.JoinHostPort(k.targetNodeAddr, k.targetNodeListenPort)
	jsonValue, marshalError := json.Marshal(map[string]any{
		"type": "set",
		"key": key,
		"value": value,
	})

	if marshalError != nil {
		return marshalError
	}

	var err error
	for {
		time.Sleep(k.heartbeatingTime)

		fmt.Printf("foo")

		k.conn, err = net.Dial("tcp", joined)
		if err != nil {
			return err
		}
		k.conn.Write(jsonValue)

		reply := make([]byte, 2024)
		n, _ := k.conn.Read(reply)
		json.Unmarshal(reply[:n], &response)

		fmt.Println(response.Ack)
		k.conn.Close()
	}
	return nil
}

func (k *KnucklesDBClient) Get(key []byte) (string, error) {
	var (
		err error
		serverResponse ServerMessages 
	)

	if okKey := isEmpty(key); okKey {
		return "", errors.New("The Key is Empty!")
	}

	jsonGetValue, marshalError := json.Marshal(map[string]any{
		"type": "get",
		"key": key,
	})

	if marshalError != nil {
		return "", marshalError
	}

	joined := net.JoinHostPort(k.targetNodeAddr, k.targetNodeListenPort)

	k.conn, err = net.Dial("tcp", joined)
	if err != nil {
		return "", err
	}

	k.conn.Write(jsonGetValue)

	reply := make([]byte, 2024)
	n, _ := k.conn.Read(reply)
	json.Unmarshal(reply[:n], &serverResponse)

	k.conn.Close()

	return serverResponse.Ack, nil
}

// TODO
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