/*	Copyright [2024] [Elia Renzoni]
*
*   Licensed under the Apache License, Version 2.0 (the "License");
*   you may not use this file except in compliance with the License.
*   You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*   Unless required by applicable law or agreed to in writing, software
*   distributed under the License is distributed on an "AS IS" BASIS,
*   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*   See the License for the specific language governing permissions and
*   limitations under the License.
*/




package consensus

import (
	"fmt"
	"knucklesdb/swim"
	"knucklesdb/wal"
	_ "sync"
	"time"
	"net"
)

type AntiEntropy struct {
	gossipProtocol *Gossip
	replicaHost, replicaPort string
	membershipList *swim.ClusterManager
	infoLogger     *wal.InfoLogger
	sleepTime      func(time.Duration)
}

func NewAntiEntropy(replicaHost, replicaPort string, gProtocol *Gossip, clusterManager *swim.ClusterManager, logger *wal.InfoLogger) *AntiEntropy {
	return &AntiEntropy{
		gossipProtocol: gProtocol,
		replicaHost: replicaHost,
		replicaPort: replicaPort,
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

		fmt.Printf("Cluster Len ----> %d", a.membershipList.GetClusterLen())
		if a.membershipList.GetClusterLen() == 2 {
			if net.JoinHostPort(a.replicaHost, a.replicaPort) != a.membershipList.GetSeedNodeInfos() {
				continue
			}
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