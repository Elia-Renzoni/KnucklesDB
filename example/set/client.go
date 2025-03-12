package main

import (
	"knucklesdb/client"
	"time"
	"fmt"
)

func main() {
	knucklesClient := client.NewClientSet("127.0.0.1:5050", 1 * time.Second)
	
	for i := 0; i < 300; i++ {
		go func() {
			err := knucklesClient.Set([]byte(fmt.Sprintf("%s%d", "/foo", i)), []byte("192.78.4.1"))
			fmt.Println(err)
		}()
	}
	for {}
}
