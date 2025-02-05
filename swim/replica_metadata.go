/*
*	This module contains the basic informations about nodes.
*	According to SWIM procotol status must be ALIVE, SUSPICIOUS and REMOVED.
*
**/

package swim

import (
	"net"
)

const (
	STATUS_ALIVE int = iota * 1
	STATUS_SUSPICIOUS
	STATUS_REMOVED
)

type Node struct {
	nodeAddress net.IP
	nodeListenPort int

	// this field will containt the status
	nodeStatus int
}

func NewNode(nodeAddress net.IP, listenPort, nodeStatus int) *Node {
	return &Node{
		nodeAddress: nodeAddress,
		nodeListenPort: listenPort,
		nodeStatus: nodeStatus,
	}
}