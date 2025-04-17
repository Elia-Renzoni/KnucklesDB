package consensus

import (
	"knucklesdb/vvector"
)

type InfectionBuffer struct {
	buffer chan Entry
	serializedEntriesToSpread []byte
	versionVectorMarshaler *vvector.VersionVectorMessage
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
			serializedEntriesToSpread = append(serializedEntriesToSpread, encodedMessage)
		}
	}
}