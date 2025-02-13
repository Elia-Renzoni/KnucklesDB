package main

import (
	"knucklesdb/client/interface"
	"fmt"
)

func main() {
	client := interface.NewClient("127.0.0.1", "6800", "127.0.0.1", "5050", 0)
	result, err := client.Get([]byte("/foo"))
	if err != nil {
		fmt.Printf("%v \n", err)
	}

	fmt.Printf("%s", result)
}