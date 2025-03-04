package main

import (
	"knucklesdb/client"
	"time"
	"fmt"
)

func main() {
	knucklesClient := client.NewClient("127.0.0.1", "5050", "127.0.0.1", "5050", 3 * time.Second)
	
	go func() {
		err := knucklesClient.Set([]byte("/foo"), []byte("192.78.4.1"))
		fmt.Println(err)
	}()
	for {}
}
