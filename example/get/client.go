package main

import (
	"knucklesdb/client"
	"fmt"
	"time"
)

func main() {
	knucklesClient := client.NewClientGet("127.0.0.1:5050")
	for {
		time.Sleep(2 * time.Second)
		result, err := knucklesClient.Get([]byte("/foo"))
		if err != nil {
			fmt.Printf("%v \n", err)
		}

		fmt.Printf("%s", result)
	}
}
