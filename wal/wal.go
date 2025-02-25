package wal

import (
	"os"
	"bufio"
)

type WAL struct {
	path        string
	writeOffset int64
	readOffset  int64

	// hash + offset
	walHash map[int32]int64
	walFile *os.File
	scanner *bufio.Scanner
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
	var err error 

	w.walFile, err = os.Open(w.path)
	if err != nil {
		return 
	}
	defer w.walFile.Close()
	w.writeOffset = w.getLatestOffset()
}

func (w *WAL) IsWALFull() bool {
	return false
}

func (w *WAL) ScanLines() (key []byte, value []byte) {
	return
}

func (w *WAL) getLatestOffset() int64 {
	w.scanner = bufio.NewScanner(w.walFile)
	offset := int64(0)

	for w.scanner.Scan() {
		offset += len(w.scanner.Bytes()) + 1
	}

	return offset
}

func (w *WAL) addOffsetEntry() {

}