package consensus

import (
	"knucklesdb/vvector"
	"sync"
	"slices"
)

type InfectionBuffer struct {
	buffer chan Entry
	serializedEntriesToSpread []byte
	versionVectorMarshaler *vvector.VersionVectorMessage
	lock sync.Mutex
}

func NewInfectionBuffer(marshaler *vvector.VersionVectorMessage) *InfectionBuffer {
	return &InfectionBuffer{
		buffer: make(chan Entry),
		serializedEntriesToSpread: make([][]byte, 0),
	}
}

func (i *InfectionBuffer) WriteInfection(entryToWrite Entry) {
	i.buffer <- entryToWrite
}

func (i *InfectionBuffer) ReadInfectionToSpread() {
	for {
		select {
		case entry := <- i.buffer:
			encodedMessage, err := i.versionVectorMarshaler.MarshalVersionVectorMessage(entry.key, entry.value, entry.version)
			if err != nil {
				// log error
			}
			i.addEntryToTheSlice(encodedMessage)
		}
	}
}

func (i *InfectionBuffer) addEntryToTheSlice(entry []byte) {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.serializedEntriesToSpread = append(i.serializedEntriesToSpread, entry)
}

func (i *InfectionBuffer) DeleteEntriesFromSlice() {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.serializedEntriesToSpread = slices.Delete(i.serializedEntriesToSpread, 1, 5)
}