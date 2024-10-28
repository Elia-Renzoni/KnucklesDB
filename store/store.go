package store

import (
	"errors"
	"sync"
	"net"
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

	_, check := address.(net.IP) 
	if !check {
		return errors.New("The Value is Not an IP Address")

	}

	_, ok := k.LRUCache[address]

	if ok {
		return errors.New("The IP Address Already Exist")
	}

	k.LRUCache[address] = values
	return
}

func (k *KnucklesDB) SetWithEndpointOnly(endpoint string, values *DBvalues) (err error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	_, ok := k.LRUCache[endpoint]
	if ok {
		return errors.New("The Endpoint Already Exist")
	}

	k.LRUCache[endpoint] = values
	return
}

func (k *KnucklesDB) SearchWithIpOnly(addres string) (values *DBvalues, err error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	_, check := addres.(net.IP)
	if !check {
		return errors.New("The Value is Not an IP Address")
	}

	_, ok := k.LRUCache[addres]
	if !ok {
		return nil, errors.New("Not Found")
	}

	values = k.LRUCache[addres]
	return
}

func (k *KnucklesDB) SearchWithEndpointOnly(enpoint string) (values *DBvalues, err error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	_, ok := k.LRUCache[enpoint]
	if !ok {
		return nil, errors.New("Not Found")
	}

	values = k.LRUCache[enpoint]
	return
}

func (k *KnucklesDB) DeleteEntry(entryID string) (err error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	_, ok := k.LRUCache[entryID]
	if !ok {
		return errors.New("Not Found")
	}

	delete(k.LRUCache, entryID)
	return
}
