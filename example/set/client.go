package main

import (
	"knucklesdb/client"
	"time"
	"fmt"
)

func main() {
	knucklesClient := client.NewClientSet("127.0.0.1:5050", 3 * time.Second)
	
	for i := 0; i < 30; i++ {
		go func() {
			err := knucklesClient.Set([]byte("/foo"), []byte("192.78.4.1"))
			fmt.Println(err)
		}()
	}
	for {}
}
