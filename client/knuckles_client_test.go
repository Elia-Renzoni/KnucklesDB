package client_test

import (
	"testing"
	"knucklesdb/client"
	"time"
)

func TestIsEmpty(t *testing.T) {
	kClient := client.NewClient("127.0.0.1", "5050", "127.0.0.1", "5050", 3 *time.Second)

	if ok := kClient.IsEmpty([]byte("foo")); ok {
		t.Fail()
	}
}