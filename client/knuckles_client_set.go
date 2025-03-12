package client

import (
	"net"
	"time"
	"fmt"
	"errors"
	"encoding/json"
)

type ClientSet struct {
	targetNodeAddress string 
	sleepInterval time.Duration
	conn net.Conn
}

type ServerMessages struct {
	Ack string `json:"ack"`
}

func NewClientSet(address string, sleepTime time.Duration) *ClientSet {
	return &ClientSet{
		targetNodeAddress: address,
		sleepInterval: sleepTime,
	}
}

func (c *ClientSet) Set(key, value []byte) error {
	var response ServerMessages

	fmt.Printf("yooo")

	if okKey := c.IsEmpty(key); okKey {
		return errors.New("The Key is Empty!")
	}

	if okValue := c.IsEmpty(value); okValue {
		return errors.New("The Value is Empty!")
	}

	jsonValue, _ := json.Marshal(map[string]any{
		"type": "set",
		"key": key,
		"value": value,
	})

	fmt.Println(string(jsonValue))

	var err error
	for {
		time.Sleep(c.sleepInterval)

		//fmt.Printf("foo")

		c.conn, err = net.Dial("tcp", c.targetNodeAddress)
		if err != nil {
			return err
		}
		c.conn.Write([]byte(jsonValue))

		reply := make([]byte, 2024)
		n, _ := c.conn.Read(reply)
		json.Unmarshal(reply[:n], &response)

		c.conn.Close()
	}
	return nil
}

func (c *ClientSet) IsEmpty(bytesToCheck []byte) bool {
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