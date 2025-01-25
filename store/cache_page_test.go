package store_test


import (
	"testing"
	"knucklesdb/store"
	"fmt"
	_"bytes"
)

func TestAddPage(t *testing.T) {
	var page = store.Palloc(0)
	page.AddPage([]byte("/foo"), []byte("127.0.0.1"), 0)
	page.AddPage([]byte("/bar"), []byte("192.80.09.12"), 1)
	page.AddPage([]byte("/qux"), []byte("192.70.34.23"), 2)
	page.AddPage([]byte("/mock"), []byte("192.30.23.56"), 5)

	err, value := page.ReadValueFromBucket([]byte("/mock"))
	fmt.Println(err)
	fmt.Printf("%s \n", string(value))
	if err != nil {
		t.Fail()
	}
}

func TestDeleteBucket(t *testing.T) {
	var page = store.Palloc(0)
	page.AddPage([]byte("/foo1"), []byte("127.0.0.1"), 0)
	page.AddPage([]byte("/foo2"), []byte("192.89.22.3"), 0)

	if ok := page.DeleteBucket([]byte("/foo2")); ok == false {
		t.Fail()
	}
}