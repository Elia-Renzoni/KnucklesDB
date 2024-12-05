package main

import (
	"knucklesdb/clock"
	"knucklesdb/detector"
	"knucklesdb/server/node"
	"knucklesdb/store"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup{}
	var detectionInsertionWaiting = make(chan struct{})

	internalClock := clock.NewLogicalClock(5, 3000)
	database := store.NewKnucklesDB()
	dbValues := store.NewDBValues()
	helper := detector.NewHelper()
	detectionTree := detector.NewDectionBST()
	failureDetector := detector.NewFailureDetector(detectionTree, helper, detectionInsertionWaiting)
	replica := node.NewReplica("127.0.0.1", "5050", internalClock, database, dbValues)

	go internalClock.IncrementLogicalClock()
	go failureDetector.FaultDetection()
	go func() {
		for {
			time.Sleep(time.Duration(5))
			wg.Add(1)

			entries := database.ReturnEntries()
			// tree creation
			for _, value := range entries {
				detectionTree.Insert(value.nodeID, value.clock)
			}

			// unlock the detection goroutine
			detectionInsertionWaiting <- struct{}{}
			wg.Done()
		}
	}()
	go helper.StartEvictionProcess()

	replica.Start()
}
