package wal_test

import (
	"testing"
	"knucklesdb/wal"
)

func WriteWALTest(t *testing.T) {
	logger := wal.NewWAL()

	logger.WriteWAL(wal.NewWALEntry(3428282828, []byte("Set"), []byte("/foo"), []byte("value")))

}