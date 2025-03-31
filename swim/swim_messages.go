/**
*	This module contains the messages that the protocol will receive
*
**/
package swim

type PiggyBackMessage struct {
	MethodType string `json:"type"`
	ParentNode string `json:"node"`
	TargetNode string `json:"target"`
	PingValue  int    `json:"ping"`
}

// this message will be send over gossip by each nodes that
// pings the cluster
type DetectionMessage struct {
	MethodType    string `json:"type"`
	Status        int    `json:"swim"`
	NodeID        string `json:"nodeID"`
	NodeListePort int    `json:"port"`
}

// this message will be received as a ACK
// if ackContent is 0 then the target node is not alive
// if ackContent is 1 then the target node i alive
type AckMessage struct {
	AckContent int `json:"ack,omitempty"`
}

type JoinMessage struct {
	MethodType string `json:"type"`
	IPAddr     string `json:"ip"`
	ListenPort string `json:"port"`
}

type MembershipListMessage struct {
	MethodType string            `json:"type"`
	List       []MembershipEntry `json:"list"`
}

type MembershipEntry struct {
	NodeAddress    string `json:"address"`
	NodeListenPort string `json:"port"`
	NodeStatus     int    `json:"status"`
}