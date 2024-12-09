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
	var wg sync.WaitGroup
	var detectionInsertionWaiting = make(chan struct{})

	internalClock := clock.NewLogicalClock(5, 3000)
	database := store.NewKnucklesDB()
	helper := detector.NewHelper(wg, database)
	detectionTree := detector.NewDetectionBST()
	failureDetector := detector.NewFailureDetector(detectionTree, helper, wg, detectionInsertionWaiting)
	replica := node.NewReplica("127.0.0.1", "5050", internalClock, database)

	go internalClock.IncrementLogicalClock()
	go func(){
		for {
			time.Sleep(time.Second*5)
			//wg.Add(1)

			entries := database.ReturnEntries()
			if len(entries) != 0 {
				wg.Add(1)
				// tree creation
				for _, value := range entries {
					detectionTree.Insert(value.NodeID, value.Clock)
				}

				// unlock the detection goroutine
				detectionInsertionWaiting <- struct{}{}
				wg.Done()
			}
		}
	}()

	go failureDetector.FaultDetection()
	go helper.StartEvictionProcess()

	replica.Start()
}
