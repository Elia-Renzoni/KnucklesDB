package main

import (
	"knucklesdb/client"
	"fmt"
	_"time"
)

func main() {
	knucklesClient := client.NewClientGet("127.0.0.1:5050")
	for {
		result, err := knucklesClient.Get([]byte("/foo6"))
		if err != nil {
			fmt.Printf("%v \n", err)
		}

		fmt.Printf("%s", result)
	}
}
