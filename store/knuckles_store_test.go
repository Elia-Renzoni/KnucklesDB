/*	Copyright [2024] [Elia Renzoni]
*
*   Licensed under the Apache License, Version 2.0 (the "License");
*   you may not use this file except in compliance with the License.
*   You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*   Unless required by applicable law or agreed to in writing, software
*   distributed under the License is distributed on an "AS IS" BASIS,
*   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*   See the License for the specific language governing permissions and
*   limitations under the License.
*/



package store_test


import (
	"testing"
	"knucklesdb/store"
)


func TestSet(t *testing.T) {
	var (
		bPool = store.NewBufferPool()
		binder = store.NewAddressBinder()
		hash = store.NewSpookyHash(1)
		kMap = store.NewKnucklesMap(bPool, binder, hash)
	)

	kMap.Set([]byte("/foo"), []byte("192.78.55.55"))
	kMap.Set([]byte("/qux"), []byte("192.245.123.60"))
	kMap.Set([]byte("/bar"), []byte("192.124.255.255"))
	kMap.Set([]byte("/mock"), []byte("192.170.89.233"))

	err1, value1 := kMap.Get([]byte("/foo"))
	err2, value2 := kMap.Get([]byte("/qux"))
	err3, value3 := kMap.Get([]byte("/bar"))
	err4, value4 := kMap.Get([]byte("/mock"))

	t.Log(string(value1))
	t.Log(string(value2))
	t.Log(string(value3))
	t.Log(string(value4))

	if err1 != nil {
		t.Fail()
	}

	if err2 != nil {
		t.Fail()
	}

	if err3 != nil {
		t.Fail()
	}

	if err4 != nil {
		t.Fail()
	}
}