package client

import (
	"net"
	"errors"
	"fmt"
	"encoding/json"
)


type ClientGet struct {
	targetNodeAddress string
	conn net.Conn
}

type ServerMessagesGet struct {
	Ack string `json:"ack"`
}

func NewClientGet(target string) *ClientGet {
	return &ClientGet{
		targetNodeAddress: target,
	}
}

func (c *ClientGet) Get(key []byte) (string, error){
	var (
		err error
		serverResponse ServerMessagesGet
	)

	if okKey := c.IsEmpty(key); okKey {
		return "", errors.New("The Key is Empty!")
	}

	jsonGetValue, marshalError := json.Marshal(map[string]any{
		"type": "get",
		"key": key,
	})

	if marshalError != nil {
		return "", marshalError
	}

	c.conn, err = net.Dial("tcp", c.targetNodeAddress)
	if err != nil {
		return "", err
	}

	c.conn.Write(jsonGetValue)

	reply := make([]byte, 2024)
	n, _ := c.conn.Read(reply)
	json.Unmarshal(reply[:n], &serverResponse)

	fmt.Println(string(reply))
	c.conn.Close()

	return serverResponse.Ack, nil
}

func (c *ClientGet) IsEmpty(bytesToCheck []byte) bool {
	var resultToReturn bool = true
	for _, b := range bytesToCheck {
		if b != 0 {
			resultToReturn = false
		} else {
			resultToReturn = true
		}
	}
	return resultToReturn
}