package clock

import (
	"time"
)

type KnucklesClock interface {
	IncrementLogicalClock()
	GetLogicalClock() int16
}

type LogicalClock struct {
	// logical clock
	knucklesLogicalClock int16

	// sleep time
	timing int

	// maximum clock value
	maxValueLimiter int16

	// bidirectional channel in which write values
	clocks chan int16
}

func NewLogicalClock(timing int, limiter int16) *LogicalClock {
	switch {
	case timing <= 0 || limiter <= 0:
		return nil
	}
	return &LogicalClock{
		knucklesLogicalClock: 0,
		timing:               timing,
		maxValueLimiter:      limiter,
		clocks:               make(chan int16),
	}
}

func (l *LogicalClock) IncrementLogicalClock() {
	for {
		time.Sleep(time.Duration(l.timing))

		l.knucklesLogicalClock++

		if l.knucklesLogicalClock >= l.maxValueLimiter {
			l.knucklesLogicalClock = 0
		}
		writeClock(l.clocks, l.knucklesLogicalClock)
	}
}

func (l *LogicalClock) GetLogicalClock() int16 {
	return <-l.clocks
}

func writeClock(clocksChan chan<- int16, logicalClock int16) {
	clocksChan <- logicalClock
}
