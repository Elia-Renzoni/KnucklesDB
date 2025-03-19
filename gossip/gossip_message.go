package gossip

type GossipMessage[T any] struct {
    message T
}

func (g *GossipMessage[T]) EncodeJSON() {

}

func (g *GossipMessage[T]) DecodeJSON() {

}