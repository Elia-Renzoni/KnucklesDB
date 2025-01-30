package main

import (
	"sync"
	"knucklesdb/detector"
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
	storeMap := store.NewKnucklesMap(bufferPool, addressBind, hashAlgorithm)
	replica := node.NewReplica(*host, *port, storeMap)

	failureDetector := detector.NewDetectorBuffer(bufferPool, wg)
	updateQueue := detector.NewSingularUpdateQueue(failureDetector)

	go failureDetector.ClockPageEviction()
	go updateQueue.UpdateQueueReader()

	replica.Start()
}
