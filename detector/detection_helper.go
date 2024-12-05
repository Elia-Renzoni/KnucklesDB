package detector

import (
	"knucklesdb/store"
	"net"
	"slices"
	"sync"
)


type Helper struct {
	nodesToEvict []NodeValues
	wg *sync.WaitGroup
}

func NewHelper(wg *sync.WaitGroup) *Helper {
	return &Helper{
		nodesToEvict: make([]NodeValues, 0),
		wg: wg,
	}
}

func (h *Helper) StartEvictionProcess() {	
	for {
		h.wg.Wait()

		for index := range h.nodesToEvict {
			node := h.nodesToEvict[index]
			switch  {
			case node.GetIpAddress() == nil:
				fallthrough
			case node.GetOptionalEndpoint() != "":
				nodeVals, _ := store.SearchWithEndpointOnly(node.GetOptionalEndpoint())
				if nodeVals.GetLogicalClock() == node.GetLogicalClock() {
					store.Eviction(nodeVals.GetOptionalEndpoint())
				}
			case node.GetIpAddress() != nil:
				fallthrough
			case node.GetOptionalEndpoint() == "":
				nodeVals, _ := store.SearchWithIpOnly(node.GetIpAddress())
				if nodeVals.GetLogicalClock() == node.GetLogicalClock() {
					store.Eviction(nodeVals.GetIpAddress())
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