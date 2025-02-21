package wal

type WALEntry struct {
	method []byte
	key    []byte
	value  []byte
	hash   int32
}

func NewWALEntry(hash int32, parameters ...[]byte) WALEntry {
	return WALEntry{
		method: parameters[0],
		key:    parameters[1],
		value:  parameters[2],
		hash:   hash,
	}
}
