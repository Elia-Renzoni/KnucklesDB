package wal

import (
	"bufio"
	"bytes"
	"os"
	"strconv"
	"strings"
	"encoding/binary"
	"fmt"
)

type WAL struct {
	path        string
	writeOffset int64
	readOffset  int64

	// hash + offset
	walHash      map[uint32]int64
	walFile      *os.File
	helperBuffer bytes.Buffer

	RecoveryChannel chan WALEntry
}

const (
	WRITE int = iota * 1
	READ
)

func NewWAL(filePath string) *WAL {
	return &WAL{
		path:            filePath,
		writeOffset:     int64(0),
		readOffset:      int64(0),
		walHash:         make(map[uint32]int64),
		RecoveryChannel: make(chan WALEntry),
	}
}

func (w *WAL) WriteWAL(toAppend WALEntry) {
	var (
		err          error
		buffer       [][]byte
		entryToWrite []byte
		bytesHash []byte = make([]byte, binary.MaxVarintLen64)
	)

	w.walFile, err = os.Open(w.path)
	if err != nil {
		return
	}
	defer w.walFile.Close()

	fmt.Printf("*****")

	entryOffset, ok := w.walHash[toAppend.Hash]
	var newLine = bytes.NewBufferString("\n")
	
	binary.PutUvarint(bytesHash, uint64(toAppend.Hash))
	buffer = [][]byte{toAppend.Method, bytesHash, toAppend.Key, toAppend.Value, newLine.Bytes()}
	entryToWrite = bytes.Join(buffer, []byte(", "))
	if ok {
		w.walFile.WriteAt(entryToWrite, entryOffset)
	} else {
		w.setWriteOffset(entryToWrite)
		w.walHash[toAppend.Hash] = w.writeOffset

		w.walFile.WriteAt(entryToWrite, w.writeOffset)
	}
}

func (w *WAL) IsWALFull() bool {
	if len(w.walHash) > 0 {
		return true
	}
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
		hash, _ := strconv.Atoi(splittedText[1])
		key := []byte(splittedText[2])
		value := []byte(splittedText[3])

		entry := NewWALEntry(uint32(hash), method, key, value)
		w.RecoveryChannel <- entry
	}

	close(w.RecoveryChannel)
}

func (w *WAL) setWriteOffset(bytesToWrite []byte) {
	w.helperBuffer.Write(bytesToWrite)

	offset := w.writeOffset
	offset += int64(w.helperBuffer.Len())
	w.writeOffset = offset

	w.helperBuffer.Reset()
}
