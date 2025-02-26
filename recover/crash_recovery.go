package recover


import (
	"knucklesdb/store"
	"kcnuklesdb/wal"
)

type Recover struct {
	dbState *store.KnucklesMap
	wal *wal.WAL
}

func NewRecover(db *store.KnucklesMap, wal *wal.WAL) *Recover {
	return &Recover{
		db: db,
		wal: wal,
	}
}

func (r *Recover) StartRecovery() {

}