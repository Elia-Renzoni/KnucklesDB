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





/*
*	This file contains a module the marshal the swim messages
**/

package swim

import (
	"encoding/json"
)

type ProtocolMarshaer struct{}

func NewProtocolMarshaler() *ProtocolMarshaer {
	return &ProtocolMarshaer{}
}

func (p *ProtocolMarshaer) MarshalPing() ([]byte, error) {
	var (
		encodedPingMessage []byte
		err                error
	)

	encodedPingMessage, err = json.Marshal(map[string]any{
		"type": "ping",
		"ping": 0,
	})

	return encodedPingMessage, err
}

// the parentNode indicate the node who can't reach the target node
func (p *ProtocolMarshaer) MarshalPiggyBack(parentNode, targetNode string) ([]byte, error) {
	var (
		encodedMessage []byte
		err            error
	)

	encodedMessage, err = json.Marshal(map[string]any{
		"type":   "piggyback",
		"node":   parentNode,
		"target": targetNode,
		"ping":   0,
	})

	return encodedMessage, err
}

func (p *ProtocolMarshaer) MarshalSWIMDetectionMessage(nodeStatus, listenPort int, nodeID string) ([]byte, error) {
	var (
		encodedSWIMMessage []byte
		err                error
	)

	encodedSWIMMessage, err = json.Marshal(map[string]any{
		"type":   "swim-update",
		"swim":   nodeStatus,
		"nodeID": nodeID,
		"port":   listenPort,
	})

	return encodedSWIMMessage, err
}

func (p *ProtocolMarshaer) MarshalAckMessage(ackResultValue int) ([]byte, error) {
	var (
		encodedAckMessage []byte
		err               error
	)

	encodedAckMessage, err = json.Marshal(map[string]any{
		"ack": ackResultValue,
	})

	return encodedAckMessage, err
}

func (p *ProtocolMarshaer) MarshalJoinMessage(host, port string) ([]byte, error) {
	var (
		encodedJoinMessage []byte
		err                error
	)

	encodedJoinMessage, err = json.Marshal(map[string]any{
		"type": "join",
		"ip":   host,
		"port": port,
	})

	return encodedJoinMessage, err
}

func (p *ProtocolMarshaer) MarshalSingleNodeUpdate(host string, port, status int, remoteAddr string) ([]byte, error) {
	var (
		encodedUpdateToSpread []byte
		err                   error
	)

	encodedUpdateToSpread, err = json.Marshal(map[string]any{
		"type": "swim-update",
		"remote_addr": remoteAddr, 
		"node": host,
		"port": port,
		"changed": status,
	})

	return encodedUpdateToSpread, err
}
