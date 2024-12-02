package main

import (
	"knucklesdb/clock"
	"knucklesdb/detector"
	"knucklesdb/server/node"
	"knucklesdb/store"
)

func main() {
	internalClock := clock.NewLogicalClock(5, 3000)
	database := store.NewKnucklesDB()
	dbValues := store.NewDBValues()
	detectionTree := detector.NewDectionBST()
	failureDetector := detector.NewFailureDetector(detectionTree)
	replica := node.NewReplica("127.0.0.1", "5050", internalClock, database, dbValues)

	go internalClock.IncrementLogicalClock()
	go failureDetector.FaultDetection()
	go func() {
		// TODO:
		// sync tree goroutine and detectio goroutine

	}()

	replica.Start()
}
