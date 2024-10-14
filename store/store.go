package store

import "sync"

type KnucklesDB struct {
	lock sync.RWMutex
	tlb  map[string]DBvalues
}

func NewKnucklesDB() *KnucklesDB {
	return &KnucklesDB{
		tlb: make(map[string]DBvalues),
	}
}

func (k *KnucklesDB) SetWithIpAddressOnly(address string) (err error) {
	return
}

func (k *KnucklesDB) SetWithEndpointOnly(endpoint string) (err error) {
	return
}

func (k *KnucklesDB) SearchWithIpOnly(addres string) (err error) {
	return
}

func (k *KnucklesDB) SearchWithEndpointOnly(enpoint string) (err error) {
	return
}

func (k *KnucklesDB) DeleteEntry(entryID string) (err error) {
	return
}
