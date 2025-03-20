/**
*	Push-Based Gossip Protocol
**/

package gossip

import (
	"time"
)

type GossipProtocol struct {
	gossipUtils *GossipUtils
	marshaler *GossipMarshaler
}

func NewGossipProtocol(gossip *GossipUtils, marshaler *GossipMarshaler) *GossipProtocol {
	return &GossipProtocol{
		gossipUtils: gossip,
		marshaler: marshaler,
	}
}

func (g *GossipProtocol) StartGossipCycle(fanout []string, messageToSend GossipMessage[string]) {
	for index := range fanout {
		g.gossipUtils.Send(index)
	}
}


func (g *GossipProtocol) SpreadMembershipList(fanout []string, clusterMembers []*Node) {
	g.marshaler.MarshalMembershipList(clusterMembers)
	for index := range fanout {
		g.gossipUtils.Send(index, clusterMembers)
	}
}

func (g *GossipProtocol) HandleGossipRequest() {

}