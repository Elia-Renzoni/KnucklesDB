/*	Copyright [2024] [Elia Renzoni]
*
*   Licensed under the Apache License, Version 2.0 (the "License");
*   you may not use this file except in compliance with the License.
*   You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*   Unless required by applicable law or agreed to in writing, software
*   distributed under the License is distributed on an "AS IS" BASIS,
*   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*   See the License for the specific language governing permissions and
*   limitations under the License.
*/



/**
*	This file contains the implementation of the WAL pattern used for acheiving strong durability
*	WAL is written only if SET and DELETE operations occurs by using the method WriteWAL.
*	The implementation contains a method for scan (reading) the WAL file and write evry line
*	to a channel used for the recovery procedure.
**/

package wal

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	_"fmt"
)

type WAL struct {
	writeOffset int64
	readOffset  int64

	// hash + offset
	walHash      map[uint32]int64
	walFile      *os.File
	helperBuffer bytes.Buffer
	logger *ErrorsLogger

	RecoveryChannel chan WALEntry
}


func NewWAL(logger *ErrorsLogger) *WAL {
	return &WAL{
		writeOffset:     int64(0),
		readOffset:      int64(0),
		walHash:         make(map[uint32]int64),
		logger:          logger,
		RecoveryChannel: make(chan WALEntry),
	}
}

/*
*	@brief this method allow the database to write the requests in the WAL.
*	@param encoded entry to write
**/
func (w *WAL) WriteWAL(toAppend WALEntry) {
	var (
		err          error
		buffer       [][]byte
		entryToWrite []byte
	)

	w.walFile, err = os.OpenFile("wal.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		w.logger.ReportError(err)
		return
	}
	defer w.walFile.Close()

	// get the offset stored in the map for the given key
	entryOffset, ok := w.walHash[toAppend.Hash]
	var newLine = bytes.NewBufferString("\n")
	
	// prepare the message encoding
	buffer = [][]byte{toAppend.Method, toAppend.Key, toAppend.Value, newLine.Bytes()}
	entryToWrite = bytes.Join(buffer, []byte(", "))
	// check if the entry is already present in the map
	if ok {
		w.walFile.WriteAt(entryToWrite, entryOffset)
	} else {
		w.setWriteOffset(entryToWrite)
		w.walHash[toAppend.Hash] = w.writeOffset

		w.walFile.WriteAt(entryToWrite, w.writeOffset)
	}
}

/*
*	@brief this method check if the WAL file is full
*	@return boolean value indicating the result of the execution
**/
func (w *WAL) IsWALFull() bool {
	_, err := os.ReadFile("wal.txt")
	if err != nil {
		return false
	}

	return true
}

/**
*	@brief this method scans the WAL file and memorize the entries to a channel
*	used for the recovery session.
*/
func (w *WAL) ScanLines() {
	var err error

	w.walFile, err = os.OpenFile("wal.txt", os.O_RDONLY, 0644)
	if err != nil {
		w.logger.ReportError(err)
		return
	}
	defer w.walFile.Close()

	scanner := bufio.NewScanner(w.walFile)
	for scanner.Scan() {
		scannedText := scanner.Text()
		splittedText := strings.Split(scannedText, ", ")

		method := []byte(splittedText[0])
		key := []byte(splittedText[1])
		value := []byte(splittedText[2])

		entry := NewWALEntry(uint32(0), method, key, value)
		w.RecoveryChannel <- entry
	}

	close(w.RecoveryChannel)
	
	errTruncation := os.Remove("wal.txt")
	if errTruncation != nil {
		w.logger.ReportError(errTruncation)
	}
}

/*
*	@brief this method set the offset for the write operation
**/
func (w *WAL) setWriteOffset(bytesToWrite []byte) {
	w.helperBuffer.Write(bytesToWrite)

	offset := w.writeOffset
	offset += int64(w.helperBuffer.Len())
	offset += int64(32)
	w.writeOffset = offset

	w.helperBuffer.Reset()
}
