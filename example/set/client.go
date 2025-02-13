package main

import (
	"knucklesdb/client"
)


func main() {
	knucklesClient := client.NewClient("127.0.0.1", "6700", "127.0.0.1", "5050", 3)
	
	go knucklesClient.Set([]byte("/foo"), []byte("192.78.4.1"))
}
