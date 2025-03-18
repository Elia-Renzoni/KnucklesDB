package gossip

import (
	"net"
	"context"
	"time"
	"knucklesdb/wal"
	"errors"
)

type GossipUtils struct {
	gossipGlobalContext context.Context
	timeoutTime time.Duration
	errorLogger *wal.ErrorsLogger
}

func NewGossipUtils(logger *wal.ErrorsLogger) *GossipUtils {
	return &GossipUtils{
		gossipGlobalContext: context.Background(),
		timeoutTime: 1000 * time.Milliseconds,
		errorLogger: logger,
	}
}

func (g *GossipUtils) Send(nodeAddress string, gossipMessage any) {
	ctx, cancel := context.WithTimeout(g.gossipGlobalContext, g.timeoutTime)
	defer cancel()
	conn, err := net.Dial("tcp", nodeAddress)
	if err != nil {
		g.errorLogger.ReportError(err)
		return
	}
	defer conn.Close()

	// TODO
	conn.Write()

	data := make([]byte, 2024)

	select {
	case <-ctx.Done(): 
		g.errorLogger.ReportError(errors.New("Gossip Send Failed due to Context Timeout"))
	default:
		count, _ conn.Read(data)
		// Unmarshal Messages
	}
}