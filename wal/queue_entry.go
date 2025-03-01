package wal

import (
	"bytes"
)

type WALEntry struct {
	Method []byte
	Key    []byte
	Value  []byte
	Hash   uint32
}

func NewWALEntry(hash uint32, parameters ...[]byte) WALEntry {
	return WALEntry{
		Method: parameters[0],
		Key:    parameters[1],
		Value:  parameters[2],
		Hash:   hash,
	}
}

func (w WALEntry) IsSet() bool {
	if bytes.ContainsAny(w.Method, "Set") {
		return true
	}
	return false
}
