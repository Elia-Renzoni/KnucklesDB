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
	"errors"
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