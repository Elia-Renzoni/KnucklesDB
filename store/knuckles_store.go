/**
*	This file contains the interface through which the server will access the data
*   The API is straightforward, as is based of only two methods: Get and Set.
*
**/
package store

import (
	"knucklesdb/consensus"
)



type KnucklesMap struct {
	// current size of the main data structure
	size uint32

	// pointer to the buffer pool
	bufferPool *BufferPool

	// pointer to the binder class
	addressTranslator *AddressBinder

	// pointer to the hash function implementation
	hasher *SpookyHash

	updateQueue *SingularUpdateQueue

	walAPI *Recover

	bufferToInfect *consensus.InfectionBuffer
}

func NewKnucklesMap(bPool *BufferPool, t *AddressBinder, h *SpookyHash, queue *SingularUpdateQueue, walAPI *Recover,
	                buffer *consensus.InfectionBuffer) *KnucklesMap {
	return &KnucklesMap{
		size:              0,
		bufferPool:        bPool,
		addressTranslator: t,
		hasher:            h,
		updateQueue: queue,
		walAPI: walAPI,
		bufferToInfect: buffer,
	}
}

/**
*	@brief Add a new key-value pair to a bucket
*	@param key
*   @param value
**/
func (k *KnucklesMap) Set(key []byte, value []byte, version int, infectionFlag bool) {
	var (
		hash   uint32
		pageID uint32
	)

	hash = k.hasher.Hash32(key)
	pageID = k.addressTranslator.TranslateHash(hash)
	k.bufferPool.WritePage(int(pageID), key, value, version)
	k.updateQueue.AddVictimPage(NewVictim(key, int(pageID)))

	// write the operation to the WAL to reach strong durability.
	k.walAPI.SetOperationWAL(hash, key, value)

	_, _, currentVersion := k.Get(key)

	// if the write is sent by a client set the entry into the 
	// infection buffer.
	if infectionFlag {
		k.bufferToInfect.WriteInfection(consensus.NewEntry(key, value, currentVersion))
	}
}

/**
*	@brief Search a value using the given key
*	@param search key
*	@return value stored in a bucket
 */
func (k *KnucklesMap) Get(key []byte) (error, []byte, int) {
	var (
		hash   uint32
		pageID uint32
	)

	hash = k.hasher.Hash32(key)
	pageID = k.addressTranslator.TranslateHash(hash)
	err, value, version := k.bufferPool.ReadPage(int(pageID), key)
	return err, value, version
}