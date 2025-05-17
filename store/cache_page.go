/**
* This file contains the Bucket data structure
* A Bucket is equivalent to a Page stored in the Buffer Pool
*
 */

package store

import (
	"bytes"
	"errors"
	"sync"
)

const (
	PAGE_SIZE uint32 = 2024

	// whithin this character is possible to separate
	// the key from the value. In this way the
	// retrevail phase can become more speedy
	// key@value
	SEPARATOR rune = '@'

	INCREMENT_VERSION int = 0
)

type Page struct {

	// the page id is mark out by the index of the array
	pageID int

	// linked list
	collisionList *CollisionBuffer

	// mutual ex.
	mutex sync.Mutex
}

type Bucket struct {
	bucketData    [PAGE_SIZE]byte
	knucklesClock int
}

type CollisionBufferNode struct {
	bucketNode Bucket

	// version vector that keeps track of the versions
	nodeVersionVector int
	next              *CollisionBufferNode
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
		bucketNode:        bucket,
		nodeVersionVector: 0,
		next:              nil,
	}
}

/**
*	@brief This Method add a new page to the collision Linked List
*   @param key
*   @param value
**/
func (p *Page) AddPage(key, value []byte, version int) {
	var (
		b                 Bucket = Bucket{}
		node, currentNode *CollisionBufferNode
	)

	p.mutex.Lock()
	defer p.mutex.Unlock()

	b.bucketData = fillBucket(key, value)

	// create the node
	node = newCollisionBufferNode(b)

	// check if the key-value pairs are already present in the in-memory data structure
	cheatNode, ok := checkDuplicateKeys(p.collisionList.head, key)
	if ok {
		if checkValuesSimilarities(b.bucketData, cheatNode.bucketNode.bucketData) {
			cheatNode.bucketNode.bucketData = b.bucketData

			if version != INCREMENT_VERSION {
				cheatNode.nodeVersionVector = version
			}
		} else {
			cheatNode.bucketNode.bucketData = b.bucketData
			// if the parameter version is equal to 0 is time to increment
			// the version, otherwise is equal to any non zero values i just
			// update the version by overwriting it.
			if version == INCREMENT_VERSION {
				cheatNode.nodeVersionVector += 1
			} else {
				cheatNode.nodeVersionVector = version
			}
		}
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
func (p *Page) ReadValueFromBucket(key []byte) (error, []byte, int) {
	var (
		node *CollisionBufferNode = p.collisionList.head
	)

	p.mutex.Lock()
	defer p.mutex.Unlock()

	for node != nil {
		nodeBucketData := node.bucketNode.bucketData
		if result := bytes.Contains(nodeBucketData[:], key); result {
			_, valueToRetrieve, _ := bytes.Cut(nodeBucketData[:], []byte("@"))
			return nil, valueToRetrieve, node.nodeVersionVector
		}
		node = node.next
	}

	return errors.New("Cache Miss"), nil, 0
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

	p.mutex.Lock()
	defer p.mutex.Unlock()

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

func checkValuesSimilarities(newPage, oldPage [PAGE_SIZE]byte) bool {
	_, newPageValue, _ := bytes.Cut(newPage[:], []byte("@"))
	_, oldPageValue, _ := bytes.Cut(oldPage[:], []byte("@"))

	if result := bytes.Compare(newPageValue, oldPageValue); result == 0 {
		return true
	}

	return false
}
