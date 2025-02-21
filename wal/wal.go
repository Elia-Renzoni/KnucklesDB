package wal

type WAL struct {
	path        string
	writeOffset int64
	readOffset  int64

	// hash + offset
	walHash map[int32]int64
}

func NewWAL(filePath string) *WAL {
	return &WAL{
		path:        filePath,
		writeOffset: int64(0),
		readOffset:  int64(0),
		walHash:     make(map[int32]int64),
	}
}

func (w *WAL) WriteWAL(keyHash int32, key, value, op []byte) {
}

func (w *WAL) IsWALFull() bool {
	return false
}

func (w *WAL) ScanLines() (key []byte, value []byte) {
	return
}
