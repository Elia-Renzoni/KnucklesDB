package main

import (
	"flag"
	"fmt"
	"knucklesdb/server/node"
	"knucklesdb/store"
	"knucklesdb/swim"
	"knucklesdb/wal"
	"strconv"
	"sync"
	"time"
)

func main() {
	host := flag.String("h", "localhost", "a string")
	port := flag.String("p", "5050", "a string")
	timeoutDuration := 7 * time.Second
	kHelperNodes := 2
	routineSchedulingTime := 7 * time.Second
	var wg sync.WaitGroup

	flag.Parse()

	walLogger := wal.NewWAL("wal/wal.txt")
	queueUpdateLogger := wal.NewLockFreeQueue(walLogger)

	joiner := swim.NewClusterManager()
	marshaler := swim.NewProtocolMarshaler()
	swimFailureDetector := swim.NewSWIMFailureDetector(joiner, marshaler, kHelperNodes, routineSchedulingTime, timeoutDuration)

	bufferPool := store.NewBufferPool()
	addressBind := store.NewAddressBinder()
	hashAlgorithm := store.NewSpookyHash(1)

	failureDetector := store.NewDetectorBuffer(bufferPool, wg)
	updateQueue := store.NewSingularUpdateQueue(failureDetector)
	recover := store.NewRecover(queueUpdateLogger, walLogger)
	storeMap := store.NewKnucklesMap(bufferPool, addressBind, hashAlgorithm, updateQueue, recover)
	replica := node.NewReplica(*host, *port, storeMap, timeoutDuration, marshaler, joiner)

	// start recovery session if needed
	if full := walLogger.IsWALFull(); full {
		recover.StartRecovery(storeMap)
	}

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

	go failureDetector.ClockPageEviction()
	go updateQueue.UpdateQueueReader()
	go queueUpdateLogger.EntryReader()

	replica.Start()
}
