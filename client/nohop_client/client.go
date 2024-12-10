package main

import (
	"time"
	"net"
	"fmt"
)


func main() {
	// every seconds send an hearthbeat to KncklesDB

	for {
		conn, err := net.Dial("tcp", "127.0.0.1:5050")
		if err != nil {
			break
		}

		time.Sleep(time.Second)
		conn.Write([]byte(`{"type": "ip","method": "set","parameter": "192.89.12.3","port": 5050}`))

		reply := make([]byte, 2024)
		n, _ := conn.Read(reply)
		
		fmt.Printf(string(reply[:n]))

		conn.Close()
	}
}