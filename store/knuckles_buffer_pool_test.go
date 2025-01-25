package store_test

import (
	"testing"
	"knucklesdb/store"
	_"fmt"
)


func TestReadPage(t *testing.T) {
	var bPool = store.NewBufferPool()
	bPool.WritePage(1902, []byte("/foo"), []byte("192.90.2.3"), 30)
	bPool.WritePage(350, []byte("/bar"), []byte("192.69.34.2"), 45)
	bPool.WritePage(350, []byte("/qux"), []byte("192.78.33.5"), 78)


	err, value := bPool.ReadPage(1902, []byte("/foo"))
	t.Log(string(value))
	if err != nil {
		t.Fail()
	}

	err1, value1 := bPool.ReadPage(350, []byte("/bar"))
	t.Log(string(value1))
	if err1 != nil {
		t.Fail()
	}


	err2, value2 := bPool.ReadPage(2500, []byte("/foo"))
	t.Log(string(value2))
	if err2 == nil {
		t.Fail()
	}

	err3, value3 := bPool.ReadPage(350, []byte("/qux"))
	t.Log(string(value3))
	if err3 != nil {
		t.Fail()
	}
}

func TestEvictPage(t *testing.T) {
	var bPool = store.NewBufferPool()

	bPool.WritePage(2500, []byte("/qux"), []byte("192.77.33.22"), 88)

	if result := bPool.EvictPage(2500, []byte("/qux")); result == false {
		t.Fail()
	}

	err, value := bPool.ReadPage(2500, []byte("/qux"))
	t.Log(string(value))
	if err == nil {
		t.Fail()
	}
}