package wal

import (
	"os"
	"bufio"
	"bytes"
	"strings"
	"strconv"
)

type WAL struct {
	path        string
	writeOffset int64
	readOffset  int64

	// hash + offset
	walHash map[int32]int64
	walFile *os.File
	helperBuffer bytes.Buffer

	recoveryChannel chan WALEntry
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
		recoveryChannel: make(chan WALEntry),
	}
}

func (w *WAL) WriteWAL(toAppend WALEntry) {
	var (
		err error
		buffer [][]byte
		entryToWrite []byte
	) 

	w.walFile, err = os.Open(w.path)
	if err != nil {
		return 
	}
	defer w.walFile.Close()

	ok, entryOffset := w.walHash(toAppend.hash)
	buffer = [][]byte(toAppend.method, toAppend.hash, toAppend.key, toAppend.value, []byte("\n"))
	entryToWrite = bytes.Join(buffer, []byte(", "))
	if ok {
		w.walFile.WriteAt(entryToWrite, entryOffset)
	} else {
		w.setWriteOffset(entryToWrite)
		w.walHash[toAppend.hash] = w.writeOffset

		w.walFile.WriteAt(entryToWrite, w.writeOffset)
	}
}

func (w *WAL) IsWALFull() bool {
	return false
}

func (w *WAL) ScanLines() {
	var err error 

	w.walFile, err = os.Open(w.path)
	if err != nil {
		return
	}
	defer w.walFile.Close()

	scanner := bufio.NewScanner(w.walFile)
	for scanner.Scan() {
		scannedText := scanner.Text()
		splittedText := strings.Split(scannedText, ", ")

		method := []byte(splittedText[0])
		hash, _ = strconv.Atoi(splittedText[1])
		key := []byte(splittedText[2])
		value := []byte(splittedText[3])

		entry := NewWALEntry(int32(hash), method, key, value)
		w.recoveryChannel <- entry
	}

	close(w.recoveryChannel)
}

func (w *WAL) setWriteOffset(bytesToWrite []byte) {
	w.helperBuffer.Write(bytesToWrite)

	offset := w.writeOffset
	offset += w.helperBuffer.Len()
	w.writeOffset = offset

	w.helperBuffer.Reset()
}