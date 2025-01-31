/**
*	This file contains the implementation of the clock page replacement algorithm.
*
**/
package store

import (
	"sync"
	"time"
)

type DetectorBuffer struct {
	buffer     map[string]*Victim
	bufferSize int
	wg         sync.WaitGroup
	//commApi *api.API
	bufferPool *BufferPool
}

func NewDetectorBuffer(bPool *BufferPool, wg sync.WaitGroup) *DetectorBuffer {
	return &DetectorBuffer{
		buffer:     make(map[string]*Victim),
		bufferSize: 0,
		wg: wg,
		bufferPool: bPool,
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
				//d.commApi.EvictPage(victim.pageID, victim.key)
				d.bufferPool.EvictPage(victim.pageID, victim.key)
			}
		}
		d.wg.Done()
	}
}
