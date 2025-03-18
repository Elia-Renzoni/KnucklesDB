package gossip

import (
	"net"
	"context"
)

type GossipUtils struct {
	gossipGlobalContext context.Context
}

func NewGossipUtils() *GossipUtils {
	return &GossipUtils{}
}

func (g *GossipUtils) Send(nodeAddress string) {

	conn, err := net.Dial("tcp", nodeAddress)
	if err != nil {
		return
	}
}