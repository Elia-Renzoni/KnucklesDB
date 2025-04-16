package consensus

import (
	"net"
	"context"
)

type Gossip struct {
	gossipConn net.Conn
}


func NewGossip() *Gossip {
	return &Gossip{}
}

func (g *Gossip) send(host, port string) {

}


