package consensus

import (
	"knucklesdb/swim"
	"knucklesdb/wal"
	"time"
)

type AntiEntropy struct {
	gossipProtocol *Gossip
	membershipList *swim.Cluster
	infoLogger     *wal.InfoLogger
	sleepTime      func(time.Duration)
}

func NewAntiEntropy(gProtocol *Gossip, clusterList *swim.Cluster, logger *wal.InfoLogger) *AntiEntropy {
	return &AntiEntropy{
		gossipProtocol: gProtocol,
		membershipList: clusterList,
		infoLogger:     logger,
		sleepTime: func(frequency time.Duration) {
			time.Sleep(frequency)
		},
	}
}

func (a *AntiEntropy) ScheduleAntiEntropy() {
	a.infoLogger.ReportInfo("Anti-Entropy Routine ON")
	for {
		a.sleepTime(1000 * time.Millisecond)
	}
}
