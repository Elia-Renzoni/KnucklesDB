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
	walAPI *Recover
}

func NewBufferPool(recoverAPI *Recover) *BufferPool {
	return &BufferPool{
		walAPI: recoverAPI,
	}
}

/**
*	@brief This method allows you to allocate a new page
*	@param index of the array
*	@param key to store
*   @param value to store
*   @param logical clock
 */
func (b *BufferPool) WritePage(pageID int, key, value []byte, clock int) {
	var page *Page = b.pages[pageID]
	if page == nil {
		page = Palloc(clock)
		page.AddPage(key, value, clock)
		b.pages[pageID] = page
	} else {
		page.AddPage(key, value, clock)
	}
}

/**
*	@brief This method allows to read the value of a key-value pair
*	@param index of the array
*	@param key to search
*   @return miss or hit
*	@return value
 */
func (b *BufferPool) ReadPage(pageID int, key []byte) (error, []byte) {
	var page *Page = b.pages[pageID]

	if page == nil {
		return errors.New("Cache Miss"), nil
	}

	err, value := page.ReadValueFromBucket(key)
	return err, value
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
	// TODO : b.walAPI.DeleteOperationWAL()
	return result
}
