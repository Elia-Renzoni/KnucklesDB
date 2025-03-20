package gossip

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