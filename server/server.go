package main

import (
	"knucklesdb/clock"
	"knucklesdb/detector"
	"knucklesdb/server/node"
	"knucklesdb/store"
	"sync"
	"time"
	"flag"
	"fmt"
)

func main() {
	host := flag.String("h", "localhost", "a string")
	port := flag.String("p", "5050", "a string")


	var wg sync.WaitGroup
	var detectionInsertionWaiting = make(chan struct{})

	internalClock := clock.NewLogicalClock(5, 3000)
	database := store.NewKnucklesDB()
	helper := detector.NewHelper(wg, database)
	detectionTree := detector.NewDetectionBST()
	failureDetector := detector.NewFailureDetector(detectionTree, helper, wg, detectionInsertionWaiting)
	replica := node.NewReplica(*host, *port, internalClock, database)

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
					// only for debug
					fmt.Printf("%v - %v", value.NodeID, value.Clock)
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
