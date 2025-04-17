/*
*	SIR Gossip Model
*
**/

package consensus

import (
	"net"
	"context"
)

type Gossip struct {
	gossipConn net.Conn
	infectionBuffer *InfectionBuffer
}


func NewGossip(buffer *InfectionBuffer) *Gossip {
	return &Gossip{
		spreadingBuffer: make(chan Entry, 5),
		infectionBuffer: buffer,
	}
}

func (g *Gossip) send(host, port string) {

}

func (g *Gossip) prepareBuffer() []byte {

}
