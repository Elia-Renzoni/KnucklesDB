package main

import (
	"time"
	"net"
	"os"
	"fmt"
)

func main() {
	// every seconds send an hearthbeat to KncklesDB
	conn, err := net.Dial("tcp", "127.0.0.1:5050")
	if err != nil {
		os.Exit(1)
	}


    byt := []byte(`{"type":"ip",
					"method":"set",
					"parameter":"192.89.12.3",
					"port": 7075}`)

	for {
		time.Sleep(time.Second)
		_, err := conn.Write(byt)
		if err != nil {
			break
		}

		reply := make([]byte, 2024)
		_, err := conn.Read(reply)

		fmt.Printf(reply)

		if err != nil {
			break
		}
		conn.Close()
	}
}