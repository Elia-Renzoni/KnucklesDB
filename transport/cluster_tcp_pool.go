package transport

import (
	"sync"
)

type ClusterTCPConnectionPool struct {
	activeTCPConncetions sync.Map
}

func NewClusterTCPConnectionPool() *ClusterTCPConnectionPool {
	return &activeTCPConncetions{

	}
}

func (c *ClusterTCPConnectionPool) IsTCPConnectionAvaliable(tcpConnection string) (bool, sync.Cond) {
	wQueue, ok := c.activeTCPConncetions.Load(tcpConnection)
	if ok {
		if wQueue.inUse {
			return false, wQueue.waitingQueue
		} else {
			return true, nil
		}
	}
	c.activeTCPConncetions.Store(tcpConn, NewWaitingBuffer())
	return true, nil
}

func (c *ClusterTCPConnectionPool) SetTCPConnection(tcpConn string) {
	c.activeTCPConncetions.Store(tcpConn, NewWaitingBuffer())
}

func (c *ClusterTCPConnectionPool) WakeUpGoroutine(tcpConnection string) {
	wQueue, _ := c.activeTCPConncetions.Load(tcpConnection)

	wQueue.isUse = true
	wQueue.waitingQueue.Signal()
}


/*
     for {
	 	if val, ok := IsTCPConnectionAvaliable("127.0.0.1:8080"); !ok {
		    val.Wait()
		} else {
			break
		}
	 }

	 defer WakeUpGoroutine("127.0.0.1:8080")
*/