/*	Copyright [2024] [Elia Renzoni]
*
*   Licensed under the Apache License, Version 2.0 (the "License");
*   you may not use this file except in compliance with the License.
*   You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*   Unless required by applicable law or agreed to in writing, software
*   distributed under the License is distributed on an "AS IS" BASIS,
*   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*   See the License for the specific language governing permissions and
*   limitations under the License.
*/



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