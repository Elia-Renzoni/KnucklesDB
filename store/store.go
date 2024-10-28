package store

import (
	"errors"
	"sync"
)

type KnucklesDB struct {
	mutex sync.Mutex
	LRUCache   map[string]*DBvalues
}

func NewKnucklesDB() *KnucklesDB {
	return &KnucklesDB{
		LRUCache: make(map[string]*DBvalues),
	}
}

func (k *KnucklesDB) SetWithIpAddressOnly(address string, values *DBvalues) (err error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	_, ok := k.LRUCache[address]

	if ok {
		return errors.New("...")
	}

	k.LRUCache[address] = values
	return
}

func (k *KnucklesDB) SetWithEndpointOnly(endpoint string, values *DBvalues) (err error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	_, ok := k.LRUCache[endpoint]
	if ok {
		return errors.New("...")
	}

	k.LRUCache[endpoint] = values
	return
}

func (k *KnucklesDB) SearchWithIpOnly(addres string) (values *DBvalues, err error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	_, ok := k.LRUCache[addres]
	if !ok {
		return nil, errors.New("...")
	}

	values = k.LRUCache[addres]
	return
}

func (k *KnucklesDB) SearchWithEndpointOnly(enpoint string) (values *DBvalues, err error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	_, ok := k.LRUCache[enpoint]
	if !ok {
		return nil, errors.New("...")
	}

	values = k.LRUCache[enpoint]
	return
}

func (k *KnucklesDB) DeleteEntry(entryID string) (err error) {
	return
}
