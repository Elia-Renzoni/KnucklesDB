/**
*	This file contains the implementation of the singular update queue pattern used in KnucklesDB
*	to delete all the mutex managers that would force a bottleneck.
*   The entries are written in an unbounded channel of 100 elements by some store functions and readed 
*	by the method EntryReader, that will read 5 entry at a time.
**/

package wal

import (
	"fmt"
)

type WALLockFreeQueue struct {
	lockFreeQueue   chan WALEntry
	cycleIterations int
	wal             *WAL
}

func NewLockFreeQueue(wal *WAL) *WALLockFreeQueue {
	return &WALLockFreeQueue{
		lockFreeQueue:   make(chan WALEntry, 100),
		cycleIterations: 5,
		wal:             wal,
	}
}

func (wl *WALLockFreeQueue) AddEntry(entry WALEntry) {
	wl.lockFreeQueue <- entry
}

func (wl *WALLockFreeQueue) EntryReader() {
	fmt.Printf("ON...")
	for {
		for i := 0; i < wl.cycleIterations; i++ {
			entry := <-wl.lockFreeQueue
			wl.wal.WriteWAL(entry)
		}
	}
}