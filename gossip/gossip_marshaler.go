package gossip

type GossipMarshaler struct {

}

func NewGossipMarshaler() *GossipMarshaler {
	return &GossipMarshaler{

	}
}

func (g *GossipMarshaler) MarshalMembershipList(clusterMembers []*Node) ([]byte, error) {
}