package main

import (
	"flag"
	"knucklesdb/server/node"
	"knucklesdb/store"
	"knucklesdb/swim"
	"knucklesdb/wal"
	"knucklesdb/consensus"
	"knucklesdb/vvector"
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

	versionVectorMarshaler := vvector.NewVersionVectorMarshaler()

	infenctionBuffer := consensus.NewInfectionBuffer(versionVectorMarshaler, errorsLogger)
	gossipAntiEntropy := consensus.NewGossip(infectionBuffer, timeoutDuration, infoLogger)

	walLogger := wal.NewWAL(errorsLogger)
	queueUpdateLogger := wal.NewLockFreeQueue(walLogger, infoLogger)

	cluster := swim.NewCluster()
	marshaler := swim.NewProtocolMarshaler()

	spreader := swim.NewDissemination(timeoutDuration, infoLogger, errorsLogger, cluster, marshaler)
	joiner := swim.NewClusterManager(cluster, errorsLogger, spreader)

	antiEntropy := consensus.NewAntiEntropy(gossipAntiEntropy, joiner, infoLogger)

	swimFailureDetector := swim.NewSWIMFailureDetector(joiner, cluster, marshaler, kHelperNodes, routineSchedulingTime, timeoutDuration, infoLogger, errorsLogger, spreader)

	bufferPool := store.NewBufferPool()
	addressBind := store.NewAddressBinder()
	hashAlgorithm := store.NewSpookyHash(1)

	failureDetector := store.NewDetectorBuffer(bufferPool, wg, infoLogger)
	updateQueue := store.NewSingularUpdateQueue(failureDetector)
	recover := store.NewRecover(queueUpdateLogger, walLogger, infoLogger)
	storeMap := store.NewKnucklesMap(bufferPool, addressBind, hashAlgorithm, updateQueue, recover)
	replica := node.NewReplica(*host, *port, storeMap, timeoutDuration, marshaler, joiner, errorsLogger, infoLogger, spreader)

	// start recovery session if needed
	if full := walLogger.IsWALFull(); full {
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
	go infenctionBuffer.ReadInfectionToSpread()
	go antiEntropy.ScheduleAntiEntropy()

	replica.Start()
}
