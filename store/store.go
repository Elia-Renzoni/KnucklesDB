package store

import (
	"errors"
	"sync"
	"net"
	"strings"
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

	if ip := net.ParseIP(address); ip == nil {
		return errors.New("Invalid IP Address")
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

	if endpointPrefix := strings.HasPrefix(endpoint, "/"); !endpointPrefix {
		return errors.New("Invalid Endpoint")
	}
	
	_, ok := k.LRUCache[endpoint]
	if ok {
		return errors.New("The Endpoint Already Exist")
	}

	k.LRUCache[endpoint] = values
	return
}

func (k *KnucklesDB) SearchWithIpOnly(address string) (values *DBvalues, err error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	if ip := net.ParseIP(address); ip == nil {
		return nil, errors.New("Invalid IP Address")
	}

	_, ok := k.LRUCache[address]
	if !ok {
		return nil, errors.New("Not Found")
	}

	values = k.LRUCache[address]
	return
}

func (k *KnucklesDB) SearchWithEndpointOnly(endpoint string) (values *DBvalues, err error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	if endpointPrefix := strings.HasPrefix(endpoint, "/"); !endpointPrefix {
		return nil, errors.New("Invalid Endpoint")
	} 

	_, ok := k.LRUCache[endpoint]
	if !ok {
		return nil, errors.New("Not Found")
	}

	values = k.LRUCache[endpoint]
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

type NodePairs struct {
	nodeID string
	clock int16
}

type entries []NodePairs

func (k *KnucklesDB) ReturnEntries() entries {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	var pairs = make(entries, 0)

	for key, value := range k.LRUCache {
		pairs = append(pairs, NodePairs{
			nodeID: key,
			clock: value.logicalClock,
		})
	}
	return pairs
}
