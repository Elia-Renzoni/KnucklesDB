/**
*	This module contains the messages that the protocol will receive
*
**/
package swim

// this message will be send over gossip by each nodes that
// pings the cluster
type DetectionMessage struct {
	status int `json:"swim"`
	nodeID string `json:"nodeID"`	
	nodeListePort int `json:"port"`	
}

// this message will be received as a ACK 
// if ackContent is 0 then the target node is not alive
// if ackContent is 1 then the target node i alive
type AckMessage struct {
	ackContent int `json:"ack,omitempty"`	
}
