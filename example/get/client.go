package main

import (
	"knucklesdb/client"
	"fmt"
)

func main() {
	knucklesClient := client.NewClient("127.0.0.1", "6800", "127.0.0.1", "5050", 0)
	result, err := knucklesClient.Get([]byte("/foo"))
	if err != nil {
		fmt.Printf("%v \n", err)
	}

	fmt.Printf("%s", result)
}