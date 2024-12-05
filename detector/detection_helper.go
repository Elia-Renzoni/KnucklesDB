package detector

import (
	"knucklesdb/store"
	"net"
)


type Helper struct {
	nodesToEvict []NodeValues
}

func NewHelper() *Helper {
	return &Helper{
		nodesToEvict: make([]NodeValues, 0)
	}
}

func (h *Helper) StartEvictionProcess() {
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
	// TODO: delete the slice
}


func (h *Helper) AddNodeToEvict(nodes ...NodeValues) {
	for _, value := range nodes {
		h.nodesToEvict = append(h.nodesToEvict, value)
	}
}