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

func (b *BufferPool) WritePage(pageID int, key, value []byte, clock int) {
	var page *Page = Palloc(clock, key, value, pageID)
	page.AddPage(key, value)
	b.pages[pageID] = page
}

func (b *BufferPool) ReadPage(pageID int, key []byte) (error, []byte) {
	var page *Page = b.pages[pageID]

	if page == nil {
		return errors.New("Cache Miss"), nil
	}

	err, value := page.ReadValueFromBucket(key)
	return err, value
}

func (b *BufferPool) EvictPage(pageID int) {
	// TODO
}