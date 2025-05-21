/*	Copyright [2024] [Elia Renzoni]
*
*   Licensed under the Apache License, Version 2.0 (the "License");
*   you may not use this file except in compliance with the License.
*   You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*   Unless required by applicable law or agreed to in writing, software
*   distributed under the License is distributed on an "AS IS" BASIS,
*   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*   See the License for the specific language governing permissions and
*   limitations under the License.
*/




/**
*	This file contains the implementation of the singular update queue pattern used in KnucklesDB
*	to delete all the mutex managers that would force a bottleneck.
*   The entries are written in an unbounded channel of 100 elements by some store functions and readed 
*	by the method EntryReader, that will read 5 entry at a time.
**/

package wal

type WALLockFreeQueue struct {
	lockFreeQueue   chan WALEntry
	cycleIterations int
	wal             *WAL
	logger *InfoLogger
}

func NewLockFreeQueue(wal *WAL, logger *InfoLogger) *WALLockFreeQueue {
	return &WALLockFreeQueue{
		lockFreeQueue:   make(chan WALEntry, 100),
		cycleIterations: 5,
		wal:             wal,
		logger: logger,
	}
}

func (wl *WALLockFreeQueue) AddEntry(entry WALEntry) {
	wl.lockFreeQueue <- entry
}

func (wl *WALLockFreeQueue) EntryReader() {
	wl.logger.ReportInfo("WAL Reader On")	
	for {
		for i := 0; i < wl.cycleIterations; i++ {
			entry := <-wl.lockFreeQueue
			wl.wal.WriteWAL(entry)
		}
	}
}