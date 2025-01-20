package store_test

import (
	"knucklesdb/store"
	"testing"
)

func TestTranslateHash(t *testing.T) {
	binder := store.NewAddressBinder()
	bucketAddress := binder.TranslateHash(3164042272)

	if bucketAddress != 2272 {
		t.Fail()
	}
}
