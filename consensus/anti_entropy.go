package consensus

import (
	"knucklesdb/swim"
	"knucklesdb/wal"
	"time"
	"strconv"
)

type AntiEntropy struct {
	gossipProtocol *Gossip
	membershipList *swim.ClusterManager
	infoLogger     *wal.InfoLogger
	sleepTime      func(time.Duration)
}

func NewAntiEntropy(gProtocol *Gossip, clusterManager *swim.ClusterManager, logger *wal.InfoLogger) *AntiEntropy {
	return &AntiEntropy{
		gossipProtocol: gProtocol,
		membershipList: clusterManager,
		infoLogger:     logger,
		sleepTime: func(frequency time.Duration) {
			time.Sleep(frequency)
		},
	}
}

/*
*	@brief this method is responsible to schedule a foregroud gossip process.
**/
func (a *AntiEntropy) ScheduleAntiEntropy() {
	a.infoLogger.ReportInfo("Anti-Entropy Routine ON")
	for {
		a.sleepTime(1000 * time.Millisecond)
		
		// if the buffer is empty makes no sense to start a gossip
		// cycle
		if ok := a.gossipProtocol.IsBufferEmpty(); ok {
			continue
		}

		// return the buffer containing all the marshaled entries.
		encodedBufferToBroadcast := a.gossipProtocol.PrepareBuffer()

		// send the pipeline containing the version to the chosen replicas
		clusterList := a.membershipList.SetFanoutList()

		for _, nodeInfos := range clusterList {
			a.gossipProtocol.Send(nodeInfos.nodeAddress)
		}		
	}
}
