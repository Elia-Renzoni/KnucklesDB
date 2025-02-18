/*
*	This module contains the basic informations about nodes.
*	According to SWIM procotol status must be ALIVE, SUSPICIOUS and REMOVED.
*
**/

package swim


const (
	STATUS_ALIVE int = iota * 1
	STATUS_SUSPICIOUS
	STATUS_REMOVED
)

type Node struct {
	nodeAddress string
	nodeListenPort int

	// this field will containt the status
	nodeStatus int
}

func NewNode(nodeAddress string, listenPort, nodeStatus int) *Node {
	return &Node{
		nodeAddress: nodeAddress,
		nodeListenPort: listenPort,
		nodeStatus: nodeStatus,
	}
}