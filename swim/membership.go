/*
*
*
*
**/

package swim

import (
	"net"
)

type ClusterManager struct {
	// this field contains a list of nodes
	// that have joined the cluster.
	clusterMetadata []*Node
	// TODO -> gossip field...
}


func NewClusterManager() *ClusterManager {
	return &ClusterManager{
		clusterMetadata: make([]*Node, 0),
	}
}

/*
*	@brief this method will be called by the new nodes to join the cluster.
**/
func (c *ClusterManager) JoinCluster(address net.IP, port int) {
	n := NewNode(address, port, STATUS_ALIVE)
	c.clusterMetadata = append(c.clusterMetadata, n)

	// TODO -> start gossip cycle
}