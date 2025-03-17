/**
*	Push-Based Gossip Protocol
**/

package gossip

type GossipProtocol struct {
	gossipUtils *GossipUtils
}

func NewGossipProtocol() *GossipProtocol {
	return &GossipProtocol{}
}

func (g *GossipProtocol) StartGossipRound() {
}

func (g *GossipProtocol) SpreadMembershipList(fanout []string) {
	for index := range fanout {
		g.gossipUtils.Send(index)
	}
}