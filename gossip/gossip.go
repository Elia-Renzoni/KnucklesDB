/**
*	Push-Based Gossip Protocol
**/

package gossip

import (
	"time"
)

type GossipProtocol struct {
	gossipUtils *GossipUtils
}

func NewGossipProtocol(gossip *GossipUtils) *GossipProtocol {
	return &GossipProtocol{
		gossipUtils: gossip,
	}
}

func (g *GossipProtocol) StartGossipCycle(fanout []string, messageToSend GossipMessage[string]) {
	for index := range fanout {
		g.gossipUtils.Send(index)
	}
}


func (g *GossipProtocol) SpreadMembershipList(fanout []string) {
	for index := range fanout {
		g.gossipUtils.Send(index)
	}
}

func (g *GossipProtocol) HandleGossipRequest() {

}