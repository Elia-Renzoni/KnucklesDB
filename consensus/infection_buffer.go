package consensus

import (
	"knucklesdb/vvector"
	"knucklesdb/wal"
	"slices"
	"sync"
)

type InfectionBuffer struct {
	buffer                    chan Entry
	serializedEntriesToSpread []vvector.VersionVectorMessage
	lock                      sync.Mutex
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
			i.addEntryToTheSlice(entry)
		}
	}
}

func (i *InfectionBuffer) addEntryToTheSlice(entry Entry) {
	i.lock.Lock()
	defer i.lock.Unlock()

	entryVersion := vvector.VersionVectorMessage{
		Key:     entry.key,
		Value:   entry.value,
		Version: entry.version,
	}

	i.serializedEntriesToSpread = append(i.serializedEntriesToSpread, entryVersion)
}

func (i *InfectionBuffer) DeleteEntriesFromSlice() {
	i.lock.Lock()
	defer i.lock.Unlock()

	for j := 0; j < i.getNIteraction(); j++ {
		i.serializedEntriesToSpread = slices.Delete(i.serializedEntriesToSpread, j, j+1)
	}
}

func (i *InfectionBuffer) getNIteraction() int {
	if len(i.serializedEntriesToSpread) >= 5 {
		return 5
	}

	return len(i.serializedEntriesToSpread)
}

func (i *InfectionBuffer) GetFirstFiveEntries() []vvector.VersionVectorMessage {
	i.lock.Lock()
	defer i.lock.Unlock()

	return i.serializedEntriesToSpread[:5]
}
