package consensus

import (
	"bytes"
	"slices"
	"sync"
	"knucklesdb/vvector"
	"knucklesdb/wal"
)

type InfectionBuffer struct {
	buffer                    chan Entry
	serializedEntriesToSpread []Entry
	lock sync.Mutex
	versionVectorMarshaler    *vvector.VersionVectorMarshaler
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
			//encodedMessage, err := i.versionVectorMarshaler.MarshalVersionVectorMessage(entry.key, entry.value, entry.version)
			if err != nil {
				i.errorlogger.ReportError(err)
				return
			}
			i.addEntryToTheSlice(encodedMessage)
		}
	}
}

func (i *InfectionBuffer) addEntryToTheSlice(entry Entry) {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.serializedEntriesToSpread = append(i.serializedEntriesToSpread, entry)
}

func (i *InfectionBuffer) DeleteEntriesFromSlice() {
	i.lock.Lock()
	defer i.lock.Unlock()

	for i := 0; i < 5; i++ {
		i.serializedEntriesToSpread = slices.Delete(i.serializedEntriesToSpread, i, i + 1)
	}
}

func (i *InfectionBuffer) GetFirstFiveEntries() []Entry {
	i.lock.Lock()
	defer i.lock.Unlock()


	return i.serializedEntriesToSpread[:5]
}