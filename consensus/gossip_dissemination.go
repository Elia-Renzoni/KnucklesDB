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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"knucklesdb/swim"
	"knucklesdb/vvector"
	"knucklesdb/wal"
	"math"
	"net"
	"slices"
	"time"

	id "github.com/google/uuid"
)

type Gossip struct {
	gossipConn      net.Conn
	replicaHost, replicaPort string
	replicaUUID     id.UUID
	infectionBuffer *InfectionBuffer
	gossipContext   context.Context
	gossipTimeout   time.Duration
	ackMessage      swim.AckMessage
	infoLogger      *wal.InfoLogger
	errorLogger     *wal.ErrorsLogger
	terminationMap  map[id.UUID]int
	logicalClock    int
}

func NewGossip(replicaHost, replicaPort string, buffer *InfectionBuffer, uuid id.UUID, timeout time.Duration, logger *wal.InfoLogger, errorLogger *wal.ErrorsLogger) *Gossip {
	return &Gossip{
		replicaHost: replicaHost,
		replicaPort: replicaPort,	
		replicaUUID:     uuid,
		infectionBuffer: buffer,
		gossipContext:   context.Background(),
		gossipTimeout:   timeout,
		infoLogger:      logger,
		errorLogger:     errorLogger,
		terminationMap:  make(map[id.UUID]int),
		logicalClock:    0,
	}
}

func (g *Gossip) Send(address string, gossipMessage []byte) {
	ctx, cancel := context.WithTimeout(g.gossipContext, g.gossipTimeout)
	defer cancel()

	conn, err := net.Dial("tcp", address)
	if err != nil {
		g.errorLogger.ReportError(err)
		return
	}
	defer conn.Close()

	conn.Write(gossipMessage)

	data := make([]byte, 2024)

	select {
	case <-ctx.Done():
		g.errorLogger.ReportError(errors.New("Gossip Send Failed due to Context Timeout"))
	default:
		count, _ := conn.Read(data)
		json.Unmarshal(data[:count], &g.ackMessage)
	}
}

func (g *Gossip) PrepareBuffer() []vvector.VersionVectorMessage {
	entries :=  g.infectionBuffer.GetFirstFiveEntries()

	// delete the first five entries form the slice
	g.infectionBuffer.DeleteEntriesFromSlice()
	return entries
}

func (g *Gossip) IsBufferEmpty() bool {
	if len(g.infectionBuffer.serializedEntriesToSpread) >= 5 {
		return true
	}
	return false
}

func (g *Gossip) MarshalPipeline(splittedBuffer []vvector.VersionVectorMessage) ([]byte, error) {
	var (
		marshaledPipeline []byte
		err               error
	)

	g.setLogicalClockForGossipSpreading()

	marshaledPipeline, err = json.Marshal(map[string]any{
		"type":  "gossip",
		"remote_addr": net.JoinHostPort(g.replicaHost, g.replicaPort),
		"uuid":  g.replicaUUID,
		"clock": g.logicalClock,
		"data":  splittedBuffer,
	})

	fmt.Println(string(marshaledPipeline))

	return marshaledPipeline, err
}

func (g *Gossip) setLogicalClockForGossipSpreading() {
	if (g.logicalClock + 1) < math.MaxInt {
		g.logicalClock += 1
	} else {
		g.logicalClock = 1
	}
}

func (g *Gossip) AddReplicaInTerminationMap(uuid id.UUID, clock int) {
	g.terminationMap[uuid] = clock
}

func (g *Gossip) SearchReplica(uuid id.UUID) (bool, int) {
	clock, ok := g.terminationMap[uuid]
	return ok, clock
}

/*
*	@brief this method check if in the received pipelined
*	there are the same hash values. If there are the same
*	keys the method perform a partial LLW between the Pipeline
*	and then between the LLW winner and the memory content.
 */
func (g *Gossip) PipelinedLLW(pipeline []vvector.VersionVectorMessage) {

	for pipelineNodeIndex := range pipeline {
		for innerNodeIndex := range pipeline {
			outerNodeKey := pipeline[pipelineNodeIndex].Key

			if bytes.Equal(outerNodeKey, pipeline[innerNodeIndex].Key) {
				// perform a local LLW operation between entries
				outerNodeVersionVector := pipeline[pipelineNodeIndex].Version
				innerNodeVersionVector := pipeline[innerNodeIndex].Version

				switch {
				case outerNodeVersionVector > innerNodeVersionVector:
					pipeline = slices.Delete(pipeline, innerNodeIndex, innerNodeIndex)
				case innerNodeVersionVector > outerNodeVersionVector:
					pipeline = slices.Delete(pipeline, pipelineNodeIndex, pipelineNodeIndex)
				}
			}
		}
	}
}