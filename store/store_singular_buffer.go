/**
*	This file contains the implementation of a singular update queue
*	for the store package. The queue represent a single point of
*	entry both for the writers and the clock algorithm
**/

package store

type StoreSingularQeueuBuffer struct {
	updatesBuffer chan Routine
	store *KnucklesMap
	bufferPool *BufferPool
}

func NewSingularQueueBuffer(db *KnucklesMap, bPool *BufferPool) *StoreSingularQeueuBuffer {
	return &StoreSingularQeueuBuffer{
		updatesBuffer: make(chan Routine),
		store: db,
		bufferPool: bPool,
	}
}

func (s *StoreSingularQeueuBuffer) AddToBuffer(r Routine) {
	s.updatesBuffer <- r
}

// this method implement the goroutine that would read the 
// values from the channel
func (s *StoreSingularQeueuBuffer) UpdateBufferReader() {
	for {
		select {
		case r := <- s.updatesBuffer:
			switch r.operation {
			case WRITE_NEW_PAGE:
				s.store.Set(r.key, r.value)
			case EVICT_PAGE:
				s.bPool.EvictPage(r.pageID, r.key)
			}
		}
	}
}