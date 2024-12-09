package main

import (
	"time"
	"net"
	"fmt"
	"encoding/json"
)

type Payload struct {
	types string `json:"type"`
	method string `json:"method"`
	parameter string `json:"parameter"`
	port int `json:"port"`
}

type ErrorResponse struct {
	rError string `json:"error"`
}

func main() {
	// every seconds send an hearthbeat to KncklesDB
	p := &Payload{
		types: "ip",
		method: "set",
		parameter: "192.89.12.3",
		port: 5056,
	}

	payload, _ := json.Marshal(p)

	for {
		conn, err := net.Dial("tcp", "127.0.0.1:5050")
		if err != nil {
			break
		}

		time.Sleep(time.Second)
		conn.Write(payload)

		reply := make([]byte, 2024)
		conn.Read(reply)

		data := &ErrorResponse{} 
		
		json.Unmarshal(reply, data)
		fmt.Println(string(data.rError))

		conn.Close()
	}
}