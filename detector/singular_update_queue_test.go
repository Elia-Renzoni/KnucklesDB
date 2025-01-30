package detector_test

import (
	"testing"
	"knucklesdb/detector"
	"time"
	"math/rand"
	"sync"
	"knucklesdb/store"
)

func TestSingularUpdateQueue(t *testing.T) {
	var wg sync.WaitGroup
	var detectorBuffer = detector.NewDetectorBuffer(store.NewBufferPool(), wg)
	var updateQueue = detector.NewSingularUpdateQueue(detectorBuffer)

	t.Run("", func(t *testing.T) {
		t.Parallel()
		rand.Seed(time.Now().UnixNano())
		
		for {
			time.Sleep(3 * time.Second)
			str := "abcdefghilmnopqrstuvz"
			key := str[0:rand.Intn(18)]
			victim := detector.NewVictim([]byte(key), rand.Intn(3001))
			updateQueue.AddVictimPage(victim)
			t.Log("Added Page")
			ok := detectorBuffer.SearchPage(victim)
			t.Log(ok)
		}
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		updateQueue.UpdateQueueReader()
	})


	t.Run("", func(t *testing.T) {
		t.Parallel()
		detectorBuffer.ClockPageEviction()
	})
} 