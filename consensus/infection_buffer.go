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
