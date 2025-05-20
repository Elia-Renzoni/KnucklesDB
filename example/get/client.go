package main

import (
	"fmt"
	"time"
	"encoding/json"
	"net"
)

type ServerMessagesGet struct {
	Ack string `json:"ack"`
}

func main() {
	for {
		/*result, err := knucklesClient.Get([]byte("/foo6"))
		if err != nil {
			fmt.Printf("%v \n", err)
		}

		time.Sleep(2 *time.Second)
		fmt.Printf("%s", result)*/

		time.Sleep(2 *time.Second)

		var (
			err error
			serverResponse ServerMessagesGet
		)

		jsonGetValue, marshalError := json.Marshal(map[string]any{
			"type": "get",
			"key": []byte("/foo6"),
		})

		if marshalError != nil {
			return
		}

		conn, err := net.Dial("tcp", "127.0.0.1:6060")
		if err != nil {
			return
		}

		conn.Write(jsonGetValue)

		reply := make([]byte, 2024)
		n, _ := conn.Read(reply)
		json.Unmarshal(reply[:n], &serverResponse)

		fmt.Println(serverResponse)

		conn.Close()

	}
}
