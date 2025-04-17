/*
*	SIR Gossip Model
* 	In this gossip implementation the status is considered infected only if there are
*	at least 5 updates. To spread the informations among the cluster the algorithm takes
*	the most recent 5 updates and after encoding them in json, spread the updates among the 
*	cluster, by doing this the algorithm remove the most recent 5 updates from his local
*	memory structure.
**/
package consensus

import (
	"net"
	"context"
)

type Gossip struct {
	gossipConn net.Conn
	infectionBuffer *InfectionBuffer
	gossipContext context.Context
}


func NewGossip(buffer *InfectionBuffer) *Gossip {
	return &Gossip{
		spreadingBuffer: make(chan Entry, 5),
		infectionBuffer: buffer,
		gossipContext: context.Background(),
	}
}

func (g *Gossip) send(host, port string) {

}

func (g *Gossip) PrepareBuffer() (splittedBuffer [][]byte) {
	splittedBuffer = g.infectionBuffer.serializedEntriesToSpread[:5]
	
	// delete the first five entries form the slice
	g.infectionBuffer.DeleteEntriesFromSlice()
	return
}

func (g *Gossip) IsBufferEmpty() bool {
	if len(g.infectionBuffer.serializedEntriesToSpread) >= 5 {
		return true
	}
	return false
}