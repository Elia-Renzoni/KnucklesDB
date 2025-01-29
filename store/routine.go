/**
*	This file contains the implementation of the elements contained
*	in the store singular update queue buffer.
*
**/


package store

type Routine struct {
	// the value could be WRITE or EVICT
	operation int
	key []byte
	value []byte
	// useful for get the eviction phase speedy
	pageID int
}

func NewRoutine(op, pageID int, key, value []byte) Routine {
	return Routine{
		operation: op,
		key: key,
		value: value,
		pageID: pageID,
	}
}