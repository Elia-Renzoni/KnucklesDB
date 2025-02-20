package wal

type WAL struct {
	path string
}

func NewWAL() *WAL {
	return &WAL{

	}
}

func (w *WAL) appendBytes() {
}

func (w *WAL) recoveryRead() {

}