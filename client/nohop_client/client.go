package main

import (
	"time"
	"net"
	"fmt"
	_"encoding/json"
	_"math/rand"
)


func main() {
	// every seconds send an hearthbeat to KncklesDB

	for {
		conn, err := net.Dial("tcp", "127.0.0.1:5050")
		if err != nil {
			break
		}

		/*rand.Seed(time.Now().UnixNano())
		randomIPAddr := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))*/
		

		time.Sleep(1 * time.Second)
		conn.Write([]byte(`{"type": "set", "key": "/foo", "value": "bar"}`))

		reply := make([]byte, 2024)
		n, _ := conn.Read(reply)
		
		fmt.Printf(string(reply[:n]))

		conn.Close()
	}
}