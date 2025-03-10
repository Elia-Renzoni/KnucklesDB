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