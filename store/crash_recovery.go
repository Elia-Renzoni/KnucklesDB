package store

import (
	"knucklesdb/wal"
)

type Recover struct {
	walAPI             *wal.WALLockFreeQueue
	walRecoveryChannel *wal.WAL
}

func NewRecover(wal *wal.WALLockFreeQueue, walChannel *wal.WAL) *Recover {
	return &Recover{
		walAPI:             wal,
		walRecoveryChannel: walChannel,
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

func (r *Recover) StartRecovery(dbState *KnucklesMap) {
	go r.walRecoveryChannel.ScanLines()

	for {
		select {
		case entryToRestore := <-r.walRecoveryChannel.RecoveryChannel:
			if entryToRestore.IsSet() {
				dbState.Set(entryToRestore.Key, entryToRestore.Value)
			}
		// the channel is closed.
		default:
			break
		}
	}
}
