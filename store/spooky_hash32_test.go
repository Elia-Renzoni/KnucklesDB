package store_test

import (
	"knucklesdb/store"
	"testing"
)

func TestHash32(t *testing.T) {
	hasher := store.NewSpookyHash(1)
	hashValue := hasher.Hash32([]byte("/myendpoint"))

	t.Log(hashValue)
	if hashValue != 104876828 {
		t.Fail()
	}
}
