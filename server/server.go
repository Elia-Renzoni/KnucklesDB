package main

import (
	"knucklesdb/clock"
	"knucklesdb/detector"
	"knucklesdb/server/node"
	"knucklesdb/store"
	"sync"
)

func main() {
	var wg sync.WaitGroup{}

	internalClock := clock.NewLogicalClock(5, 3000)
	database := store.NewKnucklesDB()
	dbValues := store.NewDBValues()
	helper := detector.NewHelper()
	detectionTree := detector.NewDectionBST()
	failureDetector := detector.NewFailureDetector(detectionTree, helper)
	replica := node.NewReplica("127.0.0.1", "5050", internalClock, database, dbValues)

	go internalClock.IncrementLogicalClock()
	go failureDetector.FaultDetection()
	go func() {
		wg.Add()
		// TODO:
		// sync tree goroutine and detectio goroutine
		wg.Done()
	}()
	go helper.StartEvictionProcess()

	replica.Start()
}
