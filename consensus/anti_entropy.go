package consensus

import (
	"fmt"
	"knucklesdb/swim"
	"knucklesdb/wal"
	_ "sync"
	"time"
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
		if ok := a.gossipProtocol.IsBufferEmpty(); !ok {
			continue
		}

		// if the cluster is empty don't spread updates
		if !a.membershipList.ClusterLen() {
			continue
		}

		a.infoLogger.ReportInfo("Able to Spread Informations")

		// return the buffer containing the first five entries
		encodedBufferToSend := a.gossipProtocol.PrepareBuffer()

		message, _ := a.gossipProtocol.MarshalPipeline(encodedBufferToSend)
		fmt.Println(string(message))

		// send the pipeline containing the version to the chosen replicas
		clusterList := a.membershipList.SetFanoutList()

		for _, nodeInfos := range clusterList {
			a.gossipProtocol.Send(nodeInfos, message)
		}
	}
}
