package clock_test

import (
	"knucklesdb/clock"
	"testing"
)

func TestNewLogicalClock(t *testing.T) {
	var (
		timing  int   = 0
		limiter int16 = -23
	)
	if instance := clock.NewLogicalClock(timing, limiter); instance != nil {
		t.Fail()
	}
}

func TestNewLogicalClock2(t *testing.T) {
	var (
		correctTming     int   = 5
		uncorrectLimiter int16 = 0
	)
	if instance2 := clock.NewLogicalClock(correctTming, uncorrectLimiter); instance2 != nil {
		t.Fail()
	}
}

func TestNewLogicalClock3(t *testing.T) {
	var (
		uncorrectTming int   = 0
		correctLimiter int16 = 12
	)
	if instance := clock.NewLogicalClock(uncorrectTming, correctLimiter); instance != nil {
		t.Fail()
	}
}

func TestGetLogicalClock(t *testing.T) {
	instance := clock.NewLogicalClock(5, 30)
	go instance.IncrementLogicalClock()

	if clock := instance.GetLogicalClock(); clock == 0 || clock > 30 {
		t.Fail()
	}
}
