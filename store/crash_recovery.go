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
*	This file contains the implementation of the singular update queue pattern used for
*	avoiding to use a mutex that control the WAL. 
*	The implementation provide methods to write the Set and the Delete operation and also
*	provide a method to start the recovery session.
*/

package store

import (
	"knucklesdb/wal"
)

type Recover struct {
	// singular update queue for communicating
	// with the WAL.
	walAPI             *wal.WALLockFreeQueue
	
	// WAL file
	walRecoveryChannel *wal.WAL

	logger *wal.InfoLogger
}

func NewRecover(wal *wal.WALLockFreeQueue, walChannel *wal.WAL, logger *wal.InfoLogger) *Recover {
	return &Recover{
		walAPI:             wal,
		walRecoveryChannel: walChannel,
		logger: logger,
	}
}

func (r *Recover) SetOperationWAL(hash uint32, key, value []byte) {
	entry := wal.NewWALEntry(hash, []byte("Set"), key, value)

	r.walAPI.AddEntry(entry)
}

func (r *Recover) DeleteOperationWAL(hash uint32, key, value []byte) {
	entry := wal.NewWALEntry(hash, []byte("Delete"), key, value)

	r.walAPI.AddEntry(entry)
}


/**
*	@brief this method starts a producer and a consumer goroutine.
*   @param instance of the actual store.
*/
func (r *Recover) StartRecovery(dbState *KnucklesMap) {
	r.logger.ReportInfo("Starting Recovery Session")
	// start the producer
	go r.walRecoveryChannel.ScanLines()

	// start the consumer that will restore the memory content after
	// a crash.
	go func() {
		for {
			select {
			case entryToRestore := <-r.walRecoveryChannel.RecoveryChannel:
				if entryToRestore.IsSet() {
					dbState.Set(entryToRestore.Key, entryToRestore.Value, 0, true)
				}
			// the channel is closed.
			default:
				break
			}
		}
	}()
	r.logger.ReportInfo("Recovery Session Ended")
}
