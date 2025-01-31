package main

import (
	"sync"
	"knucklesdb/server/node"
	"knucklesdb/store"
	"flag"
)

func main() {
	host := flag.String("h", "localhost", "a string")
	port := flag.String("p", "5050", "a string")

	var wg sync.WaitGroup

	bufferPool := store.NewBufferPool()
	addressBind := store.NewAddressBinder()
	hashAlgorithm := store.NewSpookyHash(1)

	failureDetector := store.NewDetectorBuffer(bufferPool, wg)
	updateQueue := store.NewSingularUpdateQueue(failureDetector)
	storeMap := store.NewKnucklesMap(bufferPool, addressBind, hashAlgorithm, updateQueue)
	replica := node.NewReplica(*host, *port, storeMap)

	go failureDetector.ClockPageEviction()
	go updateQueue.UpdateQueueReader()

	replica.Start()
}
