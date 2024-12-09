package detector

import (
	"knucklesdb/store"
	"net"
	"slices"
	"sync"
)


type Helper struct {
	nodesToEvict []NodeValues
	db *store.KnucklesDB
	wg sync.WaitGroup
}

func NewHelper(wg sync.WaitGroup, db *store.KnucklesDB) *Helper {
	return &Helper{
		nodesToEvict: make([]NodeValues, 0),
		db: db,
		wg: wg,
	}
}

func (h *Helper) StartEvictionProcess() {	
	for {
		h.wg.Wait()

		for index := range h.nodesToEvict {
			node := h.nodesToEvict[index]
			if ip := net.ParseIP(node.nodeId); ip != nil {
				nodeVals, _ := h.db.SearchWithIpOnly(ip.String())
				if nodeVals.GetLogicalClock() == node.logicalClock {
					h.db.Eviction(nodeVals.GetIpAddress().String())
				}
			} else {
				nodeVals, _ := h.db.SearchWithEndpointOnly(node.nodeId)
				if nodeVals.GetLogicalClock() == node.logicalClock {
					h.db.Eviction(nodeVals.GetOptionalEndpoint())
				}
			}
		}
		slices.Delete(h.nodesToEvict, 0, len(h.nodesToEvict))
	}
}


func (h *Helper) AddNodeToEvict(nodes ...NodeValues) {
	for _, value := range nodes {
		h.nodesToEvict = append(h.nodesToEvict, value)
	}
}