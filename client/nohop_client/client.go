package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"time"
)

func main() {
	// every seconds send an hearthbeat to KncklesDB
	var key string
	var value string
	numb := flag.Int("n", 0, "an int")
	flag.Parse()
	switch *numb {
	case 0:
		key = "/foo"
		value = "192.89.23.44"
	case 1:
		key = "/bar"
		value = "192.255.66.77"
	case 2:
		key = "/mock"
		value = "192.78.255.1"
	case 3:
		key = "/qux"
		value = "192.12.33.56"
	}

	fmt.Printf("%s", key)

	for {
		conn, err := net.Dial("tcp", "127.0.0.1:5050")
		if err != nil {
			break
		}

		jsonValue, _ := json.Marshal(map[string]string{
			"type":  "set",
			"key":   key,
			"value": value,
		})

		time.Sleep(1 * time.Second)
		conn.Write(jsonValue)

		reply := make([]byte, 2024)
		n, _ := conn.Read(reply)

		fmt.Printf(string(reply[:n]))

		conn.Close()
	}
}
