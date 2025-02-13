package main

import (
	"knucklesdb/client/interface"
)


func main() {
	client := interface.NewClient("127.0.0.1", "6700", "127.0.0.1", "5050", 3)
	
	go client.Set([]byte("/foo"), []byte("192.78.4.1"))
}
