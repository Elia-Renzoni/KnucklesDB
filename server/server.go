package main

import (
	"sync"
	"knucklesdb/server/node"
	"knucklesdb/store"
	"knucklesdb/swim"
	"flag"
	"time"
	"strconv"
	"fmt"
)

func main() {
	host := flag.String("h", "localhost", "a string")
	port := flag.String("p", "5050", "a string")
	timeoutDuration := 7 * time.Second
	kHelperNodes := 2
	routineSchedulingTime := 7 * time.Second
	var wg sync.WaitGroup

	flag.Parse()

	joiner := swim.NewClusterManager()
	marshaler := swim.NewProtocolMarshaler()
	swimFailureDetector := swim.NewSWIMFailureDetector(joiner, marshaler, kHelperNodes, routineSchedulingTime, timeoutDuration)
	
	// add the server to the cluster.
	correctPort, _ := strconv.Atoi(*port)
	ok, err := joiner.IsSeed(*host, correctPort)
	if err != nil {
		// TODO -> Write to WAL.
		fmt.Printf("%v", err)
	}

	if !ok {
		joiner.JoinRequest(*host, *port)
	} else {
		go swimFailureDetector.ClusterFailureDetection()
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