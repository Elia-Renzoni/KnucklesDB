package transport

import (
	"sync"
)

type WaitingBuffer struct {
	inUse bool
	waitingQueue sync.Cond
}

func NewWaitingBuffer() *WaitingBuffer {
	return &WaitingBuffer{
		inUse: false,
		waitingQueue: sync.NewCond(&sync.Mutex{}),
	}
}