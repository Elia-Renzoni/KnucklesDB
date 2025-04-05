package swim

type Cluster struct {
	clusterMetadata []*Node
}

func NewCluster() *Cluster {
	return &Cluster{
		clusterMetadata: make([]*Node, 0),
	}
}
