package store_test

import (
	"knucklesdb/store"
	"net"
	"testing"
)

func TestValues(t *testing.T) {
	values := store.NewDBValues(net.IPv4(8, 8, 8, 8), 2340, 40, "/insertion")
	if values == nil {
		t.Fail()
	}
}
