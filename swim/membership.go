package swim

import (
	"net"
)

type ClusterManager struct {
	clusterMetadata []*Node
	// TODO -> gossip field...
}


func NewClusterManager() *ClusterManager {
	return &ClusterManager{
		clusterMetadata: make([]*Node, 0),
	}
}

func (c *ClusterManager) JoinCluster(address net.IP, port int) {
	n := NewNode(address, port, STATUS_ALIVE)
	c.clusterMetadata = append(c.clusterMetadata, n)

	// TODO -> start gossip cycle
}