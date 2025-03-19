/**
*	Push-Based Gossip Protocol
**/

package gossip

import (
	"time"
)

type GossipProtocol struct {
	gossipUtils *GossipUtils
	waitTime func(time.Duration)
}

func NewGossipProtocol(gossip *GossipUtils, interval time.Duration) *GossipProtocol {
	return &GossipProtocol{
		gossipUtils: gossip,
		waitTime: func(interval) {
			time.Sleep(interval)
		},
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