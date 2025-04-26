package consensus

import (
	"bytes"
	"knucklesdb/vvector"
	"knucklesdb/wal"
	"slices"
	"sync"
)

type InfectionBuffer struct {
	buffer                    chan Entry
	serializedEntriesToSpread bytes.Buffer
	versionVectorMarshaler    *vvector.VersionVectorMarshaler
	lock                      sync.Mutex
	errorlogger               *wal.ErrorsLogger
}

func NewInfectionBuffer(marshaler *vvector.VersionVectorMarshaler, logger *wal.ErrorsLogger) *InfectionBuffer {
	return &InfectionBuffer{
		buffer:                 make(chan Entry),
		versionVectorMarshaler: marshaler,
		errorlogger:            logger,
	}
}

func (i *InfectionBuffer) WriteInfection(entryToWrite Entry) {
	i.buffer <- entryToWrite
}

func (i *InfectionBuffer) ReadInfectionToSpread() {
	for {
		select {
		case entry := <-i.buffer:
			encodedMessage, err := i.versionVectorMarshaler.MarshalVersionVectorMessage(entry.key, entry.value, entry.version)
			if err != nil {
				i.errorlogger.ReportError(err)
				return
			}
			i.addEntryToTheSlice(encodedMessage)
		}
	}
}

func (i *InfectionBuffer) addEntryToTheSlice(entry []byte) {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.serializedEntriesToSpread.Write(entry)
	i.serializedEntriesToSpread.Write([]byte{';'})
}

func (i *InfectionBuffer) DeleteEntriesFromSlice() {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.serializedEntriesToSpread = slices.Delete(i.serializedEntriesToSpread, 1, 5)
}
