package detector_test


import (
	"testing"
	"knucklesdb/detector"
	"knucklesdb/store"
	"sync"
	"time"
)

func TestClockPageEviction(t *testing.T) {
	var wg sync.WaitGroup
	var bPool = store.NewBufferPool()
	fDetector := detector.NewDetectorBuffer(bPool, wg)

	fDetector.AddToDetectorBuffer(detector.NewVictim([]byte("foo"), 2400))
	fDetector.AddToDetectorBuffer(detector.NewVictim([]byte("bar"), 3000))

	t.Run("", func(t *testing.T) {
		t.Parallel()
		for {
			time.Sleep(2 *time.Second)
			ok := fDetector.SearchPage(detector.NewVictim([]byte("foo"), 2400))
			t.Log(ok)
		}
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()
		fDetector.ClockPageEviction()
	})
}