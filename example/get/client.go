package main

import (
	"knucklesdb/client"
	"fmt"
)

func main() {
	knucklesClient := client.NewClient("127.0.0.1", "5050", "127.0.0.1", "7070", 0)
	for {
		result, err := knucklesClient.Get([]byte("/foo"))
		if err != nil {
			fmt.Printf("%v \n", err)
		}

		fmt.Printf("%s", result)
	}
}