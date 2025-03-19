package consensus

import (
	"knucklesdb/wal"
	"knucklesdb/gossip"
	"knucklesdb/swim"
)


type AntiEntropy struct {
	gossipProtocol *gossip.GossipProtocol
	membershipList *swim.ClusterManager
}

func NewAntiEntropy(gProtocol *gossip.GossipProtocol, clusterList *swim.ClusterManager) *AntiEntropy {
	return &AntiEntropy{
		gossipProtocol: gProtocol,
		membershipList: clusterList,
	}
}

func (a *AntiEntropy) ScheduleAntiEntropy() {

}