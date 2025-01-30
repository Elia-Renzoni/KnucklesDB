package main

import (
	"net"
	"time"
	"fmt"
)


func main() {

	for {
		conn, err := net.Dial("tcp", "127.0.0.1:5050")
		if err != nil {
			break
		}

		time.Sleep(3 * time.Second)
		conn.Write([]byte(`{"type": "get", "key": "/foo"}`))

		reply := make([]byte, 2024)
		n, _ := conn.Read(reply)
		
		fmt.Printf(string(reply[:n]))

		conn.Close()
	}
	

}