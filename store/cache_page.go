/**
* This file contains the Bucket data structure
* A Bucket is equivalent to a Page stored in the Buffer Pool
*
 */

package store

import (
	"bytes"
	"errors"
)

const (
	PAGE_SIZE uint32 = 2024

	// whithin this character is possible to separate
	// the key from the value. In this way the
	// retrevail phase can become more speedy
	// key@value
	SEPARATOR rune = '@'
)

type Page struct {

	// the page id is mark out by the index of the array
	pageID int

	// logical clock increased by the clock module
	knucklesClock int

	collisionList []Bucket
}

type Bucket struct {
	bucketData [PAGE_SIZE]byte
}

func Palloc(logicalClock, bucketID int) *Page {
	return &Page{
		pageID:        bucketID,
		knucklesClock: logicalClock,
		collisionList: make([]Bucket, 0),
	}
}

func (p *Page) AddPage(key, value []byte) {
	var b Bucket = Bucket{}
	b.bucketData = fillBucket(key, value)

	p.collisionList = append(p.collisionList, b)
}

func (p *Page) ReadValueFromBucket(key []byte) (error, []byte) {
	for index := range p.collisionList {
		bucket := p.collisionList[index]
		if result := bytes.Contains(bucket.bucketData[:], key); result == true {
			_, valueToRetrieve, _ := bytes.Cut(bucket.bucketData[:], []byte("@"))
			return nil, valueToRetrieve
		}
	}

	return errors.New("Cache Miss"), nil
}

/**
*	This function fill the bucket with the datas showed as parameters
*   @param key
*	@param value
*   @return filled bucket
 */
func fillBucket(key, value []byte) (preBucket [PAGE_SIZE]byte) {
	var buffer = bytes.Buffer{}

	buffer.Write(key)
	buffer.WriteRune(SEPARATOR)
	buffer.Write(value)

	copy(preBucket[:], buffer.Bytes())
	return
}
