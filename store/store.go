package store

import (
	"errors"
	"sync"
	"net"
	"strings"
	"fmt"
)

type KnucklesDB struct {
	mutex sync.Mutex
	cache   map[string]*DBvalues
}


func NewKnucklesDB() *KnucklesDB {
	return &KnucklesDB{
		cache: make(map[string]*DBvalues),
	}
}

func (k *KnucklesDB) SetWithIpAddressOnly(address string, values *DBvalues) (err error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	if ip := net.ParseIP(address); ip == nil {
		return errors.New("Invalid IP Address")
	}

	_, ok := k.cache[address]

	if ok {
		delete(k.cache, address)
	}

	k.cache[address] = values

	// only for debug
	fmt.Printf("Set IP\n")
	return
}

func (k *KnucklesDB) SetWithEndpointOnly(endpoint string, values *DBvalues) (err error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	if endpointPrefix := strings.HasPrefix(endpoint, "/"); !endpointPrefix {
		return errors.New("Invalid Endpoint")
	}
	
	_, ok := k.cache[endpoint]
	if ok {
		delete(k.cache, endpoint)
	}

	k.cache[endpoint] = values

	// only for debug
	fmt.Printf("Set Endpoint\n")
	return
}

func (k *KnucklesDB) SearchWithIpOnly(address string) (values *DBvalues, err error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	if ip := net.ParseIP(address); ip == nil {
		return nil, errors.New("Invalid IP Address")
	}

	_, ok := k.cache[address]
	if !ok {
		return nil, errors.New("Not Found")
	}

	values = k.cache[address]
	
	// only for debug
	fmt.Printf("Get IP\n")
	return
}

func (k *KnucklesDB) SearchWithEndpointOnly(endpoint string) (values *DBvalues, err error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	if endpointPrefix := strings.HasPrefix(endpoint, "/"); !endpointPrefix {
		return nil, errors.New("Invalid Endpoint")
	} 

	_, ok := k.cache[endpoint]
	if !ok {
		return nil, errors.New("Not Found")
	}

	values = k.cache[endpoint]

	// only for debug
	fmt.Printf("Get Endpoint\n")
	return
}

func (k *KnucklesDB) Eviction(entryID string) (err error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	_, ok := k.cache[entryID]
	if !ok {
		return errors.New("Not Found")
	}

	delete(k.cache, entryID)
	return
}

type NodePairs struct {
	NodeID string
	Clock int16
}

type entries []NodePairs

func (k *KnucklesDB) ReturnEntries() entries {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	var pairs = make(entries, 0)

	for key, value := range k.cache {
		pairs = append(pairs, NodePairs{
			NodeID: key,
			Clock: value.logicalClock,
		})
	}
	return pairs
}
