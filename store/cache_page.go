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

	// linked list
	collisionList *CollisionBuffer
}

type Bucket struct {
	bucketData    [PAGE_SIZE]byte
	knucklesClock int
}

type CollisionBufferNode struct {
	bucketNode Bucket
	next       *CollisionBufferNode
}

// Buffer for the collisions that may occur
type CollisionBuffer struct {
	head *CollisionBufferNode
}

func Palloc(bucketID int) *Page {
	return &Page{
		pageID:        bucketID,
		collisionList: newCollisionBuffer(),
	}
}

func newCollisionBuffer() *CollisionBuffer {
	return &CollisionBuffer{
		head: nil,
	}
}

func newCollisionBufferNode(bucket Bucket) *CollisionBufferNode {
	return &CollisionBufferNode{
		bucketNode: bucket,
		next:       nil,
	}
}

/**
*	@brief This Method add a new page to the collision Linked List
*   @param key
*   @param value
**/
func (p *Page) AddPage(key, value []byte, logicalClock int) {
	var (
		b                 Bucket = Bucket{}
		node, currentNode *CollisionBufferNode
	)

	b.bucketData = fillBucket(key, value)
	b.knucklesClock = logicalClock

	node = newCollisionBufferNode(b)

	cheatNode, ok := checkDuplicateKeys(p.collisionList.head, key)
	if ok {
		cheatNode.bucketNode.bucketData = b.bucketData
		cheatNode.bucketNode.knucklesClock = logicalClock
	} else {
		if p.collisionList.head == nil {
			p.collisionList.head = node
		} else {
			currentNode = p.collisionList.head
			for currentNode.next != nil {
				currentNode = currentNode.next
			}

			currentNode.next = node
		}
	}
}

/**
*	@brief This Method fetch the given data in the collision list
*	@param key
*	@return error value indicating the result of the operation
*	@return value to return <key, value>
**/
func (p *Page) ReadValueFromBucket(key []byte) (error, []byte) {
	var (
		node *CollisionBufferNode = p.collisionList.head
	)

	for node != nil {
		nodeBucketData := node.bucketNode.bucketData
		if result := bytes.Contains(nodeBucketData[:], key); result {
			_, valueToRetrieve, _ := bytes.Cut(nodeBucketData[:], []byte("@"))
			return nil, valueToRetrieve
		}
		node = node.next
	}

	return errors.New("Cache Miss"), nil
}

/**
*	@brief This Method search and delete the given bucket from the collision list
*   @param key to search in list
*   @return result of the operation
**/
func (p *Page) DeleteBucket(key []byte) bool {
	var (
		node         *CollisionBufferNode = p.collisionList.head
		previousNode *CollisionBufferNode
	)

	for node != nil {
		nodeBucketData := node.bucketNode.bucketData
		if result := bytes.Contains(nodeBucketData[:], key); result {
			if previousNode == nil {
				p.collisionList.head = node.next
			} else {
				previousNode.next = node.next
			}
			return true
		}
		previousNode = node
		node = node.next
	}

	return false
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

/**
*
*	@param head of the collision linked list
*	@param key to search
*	@return pointer to node with the duplicate key
*	@return result of the search operation
 */
func checkDuplicateKeys(head *CollisionBufferNode, key []byte) (*CollisionBufferNode, bool) {
	var node *CollisionBufferNode = head

	if node == nil {
		return nil, false
	} else {
		for node != nil {
			if result := bytes.Contains(node.bucketNode.bucketData[:], key); result {
				return node, result
			}
			node = node.next
		}
	}

	return nil, false
}
