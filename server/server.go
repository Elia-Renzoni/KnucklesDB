package main

import (
	"sync"
	"knucklesdb/server/node"
	"knucklesdb/store"
	"knucklesdb/swim"
	"flag"
	"time"
)

func main() {
	host := flag.String("h", "localhost", "a string")
	port := flag.String("p", "5050", "a string")
	timeoutDuration := 10 * time.Second
	var wg sync.WaitGroup

	flag.Parse()

	joiner := swim.NewClusterManager()
	marshaler := swim.NewProtocolMarshaler()
	
	// add the server to the cluster.
	if ok := joiner.isSeed(); !ok {
		joiner.JoinRequest(*host, *port)
	}

	bufferPool := store.NewBufferPool()
	addressBind := store.NewAddressBinder()
	hashAlgorithm := store.NewSpookyHash(1)

	failureDetector := store.NewDetectorBuffer(bufferPool, wg)
	updateQueue := store.NewSingularUpdateQueue(failureDetector)
	storeMap := store.NewKnucklesMap(bufferPool, addressBind, hashAlgorithm, updateQueue)
	replica := node.NewReplica(*host, *port, storeMap, timeoutDuration, marshaler, joiner)

	go failureDetector.ClockPageEviction()
	go updateQueue.UpdateQueueReader()

	replica.Start()
}
