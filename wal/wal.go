package wal

import (
	"os"
)

type WAL struct {
	path        string
	writeOffset int64
	readOffset  int64

	// hash + offset
	walHash map[int32]int64
	walFile *os.File
}

const (
	WRITE int = iota * 1
	READ
)

func NewWAL(filePath string) *WAL {
	return &WAL{
		path:        filePath,
		writeOffset: int64(0),
		readOffset:  int64(0),
		walHash:     make(map[int32]int64),
	}
}

func (w *WAL) WriteWAL(toAppend WALEntry) {
	w.walFile = os.Open(w.path)
	w.writeOffset = w.getLatestOffset()
}

func (w *WAL) IsWALFull() bool {
	return false
}

func (w *WAL) ScanLines() (key []byte, value []byte) {
	return
}

func (w *WAL) getLatestOffset() int64 {
}

func (w *WAL) addOffsetEntry() {

}