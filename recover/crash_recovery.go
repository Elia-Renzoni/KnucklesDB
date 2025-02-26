package recover

import (
	"knucklesdb/store"
	"kcnuklesdb/wal"
)

type Recover struct {
	dbState *store.KnucklesMap
	walAPI *wal.WALLockFreeQueue
}

func NewRecover(db *store.KnucklesMap, wal *wal.WALLockFreeQueue) *Recover {
	return &Recover{
		db: db,
		wal: wal,
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
	// TODO
}