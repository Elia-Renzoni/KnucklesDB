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




/**
*	this file contains the mplementation of the buffer pool,
*    there are 3000 pages inside the buffer. Each page is uniquely indicated by the index of the array.
*	Each page contains a list of buckets
*
**/

package store


import (
	"errors"
)

const (
	MAX_PAGES int = 3000
)

type BufferPool struct {
	pages [MAX_PAGES]*Page
}

func NewBufferPool() *BufferPool {
	return &BufferPool{}
}

/**
*	@brief This method allows you to allocate a new page
*	@param index of the array
*	@param key to store
*   @param value to store
*   @param logical clock
 */
func (b *BufferPool) WritePage(pageID int, key, value []byte, version int) {
	var page *Page = b.pages[pageID]
	if page == nil {
		page = Palloc(pageID)
		page.AddPage(key, value, version)
		b.pages[pageID] = page
	} else {
		page.AddPage(key, value, version)
	}
}

/**
*	@brief This method allows to read the value of a key-value pair
*	@param index of the array
*	@param key to search
*   @return miss or hit
*	@return value
 */
func (b *BufferPool) ReadPage(pageID int, key []byte) (error, []byte, int) {
	var page *Page = b.pages[pageID]

	if page == nil {
		return errors.New("Cache Miss"), nil, 0
	}

	err, value, version := page.ReadValueFromBucket(key)
	return err, value, version
}

/**
*	@brief This method is called only by the paginator to evict pages
*	@param index of the array
*	@param key to search
*	@return result of the op.
 */
func (b *BufferPool) EvictPage(pageID int, key []byte) bool {
	var (
		page   *Page = b.pages[pageID]
		result bool
	)
	result = page.DeleteBucket(key)
	return result
}
