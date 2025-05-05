package main

import (
	"flag"
	"knucklesdb/consensus"
	"knucklesdb/server/node"
	"knucklesdb/store"
	"knucklesdb/swim"
	"knucklesdb/vvector"
	"knucklesdb/wal"
	"strconv"
	"sync"
	"time"

	id "github.com/google/uuid"
)

func main() {
	host := flag.String("h", "localhost", "a string")
	port := flag.String("p", "5050", "a string")
	timeoutDuration := 7 * time.Second
	kHelperNodes := 2
	routineSchedulingTime := 7 * time.Second
	replicaUUID := id.New()

	var (
		wg sync.WaitGroup
		syncJoin = &sync.WaitGroup{}
	)

	syncJoin.Add(1)
	flag.Parse()

	
	errorsLogger := wal.NewErrorsLogger()
	infoLogger := wal.NewInfoLogger()

	versionVectorMarshaler := vvector.NewVersionVectorMarshaler()

	infectionBuffer := consensus.NewInfectionBuffer(versionVectorMarshaler, errorsLogger)
	gossipAntiEntropy := consensus.NewGossip(infectionBuffer, replicaUUID, timeoutDuration, infoLogger, errorsLogger)

	versioningUtils := vvector.NewDataVersioning()

	walLogger := wal.NewWAL(errorsLogger)
	queueUpdateLogger := wal.NewLockFreeQueue(walLogger, infoLogger)

	cluster := swim.NewCluster()
	marshaler := swim.NewProtocolMarshaler()

	spreader := swim.NewDissemination(*host, *port, timeoutDuration, infoLogger, errorsLogger, cluster, marshaler)
	joiner := swim.NewClusterManager(syncJoin, cluster, errorsLogger, spreader)

	antiEntropy := consensus.NewAntiEntropy(gossipAntiEntropy, joiner, infoLogger)

	swimFailureDetector := swim.NewSWIMFailureDetector(joiner, cluster, marshaler, kHelperNodes, routineSchedulingTime, timeoutDuration, infoLogger, errorsLogger, spreader)

	bufferPool := store.NewBufferPool()
	addressBind := store.NewAddressBinder()
	hashAlgorithm := store.NewSpookyHash(1)

	failureDetector := store.NewDetectorBuffer(bufferPool, wg, infoLogger)
	updateQueue := store.NewSingularUpdateQueue(failureDetector)
	recover := store.NewRecover(queueUpdateLogger, walLogger, infoLogger)
	storeMap := store.NewKnucklesMap(bufferPool, addressBind, hashAlgorithm, updateQueue, recover, infectionBuffer)
	replica := node.NewReplica(*host, *port, replicaUUID, storeMap, timeoutDuration, marshaler, joiner, errorsLogger, infoLogger, spreader, gossipAntiEntropy, versioningUtils, syncJoin)

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
		go joiner.JoinRequest(*host, *port)
	}

	go swimFailureDetector.ClusterFailureDetection()
	go failureDetector.ClockPageEviction()
	go updateQueue.UpdateQueueReader()
	go queueUpdateLogger.EntryReader()
	go infectionBuffer.ReadInfectionToSpread()
	go antiEntropy.ScheduleAntiEntropy()

	replica.Start()
}
