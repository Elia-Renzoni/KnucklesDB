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
	"time"
	"knucklesdb/swim"
	"knucklesdb/wal"
	"fmt"
	"bytes"
	"encoding/json"
	"slices"
	"errors"
	"knucklesdb/vvector"
)

type Gossip struct {
	gossipConn net.Conn
	infectionBuffer *InfectionBuffer
	gossipContext context.Context
	gossipTimeout time.Duration
	ackMessage swim.AckMessage
	infoLogger *wal.InfoLogger
	errorLogger *wal.ErrorsLogger
}


func NewGossip(buffer *InfectionBuffer, timeout time.Duration, logger *wal.InfoLogger, errorLogger *wal.ErrorsLogger) *Gossip {
	return &Gossip{
		infectionBuffer: buffer,
		gossipContext: context.Background(),
		gossipTimeout: timeout,
		infoLogger: logger,
		errorLogger: errorLogger,
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

		g.infoLogger.ReportInfo(fmt.Sprintf("Ack Message: %d", g.ackMessage.AckContent))
	}
}

func (g *Gossip) PrepareBuffer() []string {
	var bufferContainingInfectionsToSend = make([]string, 0)
	
	for i := 0; i < 5; i++ {
		str, _ := g.infectionBuffer.serializedEntriesToSpread.ReadString(';')
		bufferContainingInfectionsToSend = append(bufferContainingInfectionsToSend, str)
	}
	
	// delete the first five entries form the slice
	g.infectionBuffer.DeleteEntriesFromSlice()
	return
}

func (g *Gossip) IsBufferEmpty() bool {
	if g.infectionBuffer.serializedEntriesToSpread.Len() >= 5 {
		return true
	}
	return false
}

func (g *Gossip) MarshalPipeline(splittedBuffer [][]byte) ([]byte, error) {
	var (
		marshaledPipeline []byte
		err error
	)

	marshaledPipeline, err = json.Marshal(map[string]any{
		"type": "gossip",
		"data": splittedBuffer,
	})

	return marshaledPipeline, err
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