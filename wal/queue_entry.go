package wal

import (
	"bytes"
)

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

func (w WALEntry) IsSet() bool {
	if bytes.ContainsAny(w.method, "Set") {
		return true
	}
	return false
}
