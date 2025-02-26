package recover

import (
	"knucklesdb/store"
	"kcnuklesdb/wal"
)

type Recover struct {
	dbState *store.KnucklesMap
	walAPI *wal.WALLockFreeQueue
	walRecoveryChannel *wal.WAL
}

func NewRecover(db *store.KnucklesMap, wal *wal.WALLockFreeQueue, walChannel *WAL) *Recover {
	return &Recover{
		db: db,
		wal: wal,
		walRecoveryChannel: walChannel,
	}
}

func (r *Recover) SetOperationWAL(hash int32, key, value []byte) {
	entry := wal.NewWALEntry(hash, []byte("Set"), key, value)

	r.walAPI.AddEntry(entry)
}

func (r *Recover) DeleteOperationWAL(hash int32, key, value []byte) {
	entry := wal.NewWALEntry(hash, []byte("Delete"), key, value)

	r.walAPI.AddEntry(entry)
}

func (r *Recover) StartRecovery() {
	go r.walRecoveryChannel.ScanLines()

	for {
		select {
		case entryToRestore := <- r.walRecoveryChannel:
			if entryToRestore.IsSet() {
				r.dbState.Set(entryToRestore.key, entryToRestore.value)
			}
		// the channel is closed.
		default:
			break
		}
	}
}