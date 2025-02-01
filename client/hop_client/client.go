package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"time"
)

func main() {
	var key string
	numb := flag.Int("n", 0, "an int")
	flag.Parse()
	switch *numb {
	case 0:
		key = "/foo"
	case 1:
		key = "/bar"
	case 2:
		key = "/mock"
	case 3:
		key = "/qux"
	}
	fmt.Printf("%s", key)
	for {
		conn, err := net.Dial("tcp", "127.0.0.1:5050")
		if err != nil {
			break
		}

		jsonValue, _ := json.Marshal(map[string]string{
			"type": "get",
			"key":  key,
		})
		time.Sleep(3 * time.Second)
		conn.Write(jsonValue)

		reply := make([]byte, 2024)
		n, _ := conn.Read(reply)

		fmt.Printf(string(reply[:n]))

		conn.Close()
	}
}
