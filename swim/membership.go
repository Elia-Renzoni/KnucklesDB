package swim

type ClusterManager struct {
	clusterMetadata []*Node
	// gossip...
}


func NewClusterManager() *ClusterManager {
	return &ClusterManager{
		clusterMetadata: make([]*Node, 0),
	}
}

func (c *ClusterManager) JoinCluster() {

}