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

	page.AddPage([]byte("/todel"), []byte("192.66.255.255"), 0)
	page.AddPage([]byte("/todel"), []byte("192.66.245.255"), 0)

	err1, value := page.ReadValueFromBucket([]byte("/todel"))
	t.Log(string(value))
	if err1 != nil {
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