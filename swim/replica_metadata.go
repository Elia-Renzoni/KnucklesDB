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
*	This module contains the basic informations about nodes.
*	According to SWIM procotol status must be ALIVE, SUSPICIOUS and REMOVED.
*
**/

package swim


const (
	STATUS_ALIVE int = iota * 1
	STATUS_SUSPICIOUS
	STATUS_REMOVED
)

type Node struct {
	nodeAddress string
	nodeListenPort int

	// this field will containt the status
	nodeStatus int
}

func NewNode(nodeAddress string, listenPort, nodeStatus int) *Node {
	return &Node{
		nodeAddress: nodeAddress,
		nodeListenPort: listenPort,
		nodeStatus: nodeStatus,
	}
}