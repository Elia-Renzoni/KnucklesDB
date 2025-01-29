/**
*	This file contains the implementation of the clock page replacement algorithm.
*
**/
package detector

import (
	"sync"
	"time"
	"knucklesdb/store"
)

type DetectorBuffer struct {
	buffer     map[string]*Victim
	bufferSize int
	wg         *sync.WaitGroup
	storeSingularUpdater *store.StoreSingularQeueuBuffer
}

func NewDetectorBuffer(storeQueue *store.StoreSingularQeueuBuffer) *DetectorBuffer {
	return &DetectorBuffer{
		buffer:     make(map[string]*Victim),
		bufferSize: 0,
		storeSingularUpdater: storeQueue,
	}
}

func (d *DetectorBuffer) AddToDetectorBuffer(v *Victim) {
	d.buffer[string(v.key)] = v
	d.bufferSize = len(d.buffer)
}

func (d *DetectorBuffer) SearchPage(v *Victim) bool {
	_, ok := d.buffer[string(v.key)]
	if ok {
		return ok
	}
	return false
}

func (d *DetectorBuffer) UpdateVictimEpoch(v *Victim) {
	d.buffer[string(v.key)] = v
}

func (d *DetectorBuffer) ClockPageEviction() {
	for {
		time.Sleep(3 * time.Second)

		// if the buffer is empty do nothing
		if d.bufferSize == 0 {
			continue
		}

		d.wg.Add(1)

		// loop over the buffer
		for _, victim := range d.buffer {
			if victim.epoch {
				victim.epoch = false
			} else {
				// if the page is false then i can remove it
				routine := store.NewRoutine(
					store.EVICT_PAGE,
					victim.key,
					[]byte(""),
					victim.pageID,
				)

				d.storeSingularUpdater.AddToBuffer(routine)
			}
		}
		d.wg.Done()
	}
}
