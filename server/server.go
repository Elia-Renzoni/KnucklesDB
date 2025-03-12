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

	errorsLogger := wal.NewErrorsLogger()
	infoLogger := wal.NewInfoLogger()

	walLogger := wal.NewWAL(errorsLogger)
	queueUpdateLogger := wal.NewLockFreeQueue(walLogger, infoLogger)

	joiner := swim.NewClusterManager(errorsLogger)
	marshaler := swim.NewProtocolMarshaler()
	swimFailureDetector := swim.NewSWIMFailureDetector(joiner, marshaler, kHelperNodes, routineSchedulingTime, timeoutDuration, infoLogger)

	bufferPool := store.NewBufferPool()
	addressBind := store.NewAddressBinder()
	hashAlgorithm := store.NewSpookyHash(1)

	failureDetector := store.NewDetectorBuffer(bufferPool, wg)
	updateQueue := store.NewSingularUpdateQueue(failureDetector)
	recover := store.NewRecover(queueUpdateLogger, walLogger, infoLogger)
	storeMap := store.NewKnucklesMap(bufferPool, addressBind, hashAlgorithm, updateQueue, recover)
	replica := node.NewReplica(*host, *port, storeMap, timeoutDuration, marshaler, joiner, errorsLogger, infoLogger)

	// start recovery session if needed
	if full := walLogger.IsWALFull(); full {
		fmt.Printf("here")
		recover.StartRecovery(storeMap)
	}

	correctPort, _ := strconv.Atoi(*port)
	ok, err := joiner.IsSeed(*host, correctPort)
	if err != nil {
		errorsLogger.ReportError(err)
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