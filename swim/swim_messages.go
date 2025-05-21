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
	SenderAddr string `json:"remote_addr"`
	List       []MembershipEntry `json:"list"`
}

type MembershipEntry struct {
	NodeAddress    string `json:"node"`
	NodeListenPort string `json:"port"`
	NodeStatus     int    `json:"status"`
}

type SWIMUpdateMessage struct {
	MethodType     string `json:"type"`
	SenderAddr string `json:"remote_addr"`
	NodeAddress    string `json:"node"`
	NodeListenPort int `json:"port"`
	NodeStatus     int `json:"changed"`
}
