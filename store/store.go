package store

import (
	"errors"
	"sync"
)

type KnucklesDB struct {
	mutex sync.Mutex
	tlb   map[string]*DBvalues
}

func NewKnucklesDB() *KnucklesDB {
	return &KnucklesDB{
		tlb: make(map[string]*DBvalues),
	}
}

func (k *KnucklesDB) SetWithIpAddressOnly(address string, values *DBvalues) (err error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	_, ok := k.tlb[address]

	if ok {
		return errors.New("...")
	}

	k.tlb[address] = values
	return
}

func (k *KnucklesDB) SetWithEndpointOnly(endpoint string, values *DBvalues) (err error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	_, ok := k.tlb[endpoint]
	if ok {
		return errors.New("...")
	}

	k.tlb[endpoint] = values
	return
}

func (k *KnucklesDB) SearchWithIpOnly(addres string) (values *DBvalues, err error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	_, ok := k.tlb[addres]
	if !ok {
		return nil, errors.New("...")
	}

	values = k.tlb[addres]
	return
}

func (k *KnucklesDB) SearchWithEndpointOnly(enpoint string) (values *DBvalues, err error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	_, ok := k.tlb[enpoint]
	if !ok {
		return nil, errors.New("...")
	}

	values = k.tlb[enpoint]
	return
}

func (k *KnucklesDB) DeleteEntry(entryID string) (err error) {
	return
}
