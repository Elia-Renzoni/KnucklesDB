package swim

import (
	"net"
	"wal"
)

type Dissemination struct {
	conn net.Conn
	clusterNodes []MembershipEntry
	logger *wal.InfoLogger
}

type MembershipEntry struct {
	NodeAddress string `json:"address"`
	NodeListenPort string `json:"port"`
	NodeStatus int `json:"status"`
}

func NewDissemination() *Dissemination {
	return &Dissemination{}
}

func (d *Dissemination) SpreadMembershipList(membershipList []*Node, fanoutList []string) {
	for index := range fanoutList {
		encodeClusterMetadata, err := d.marshalMembershipList(membershipList)
		d.send(fanoutList[index], encodeClusterMetadata)
	}
}

func (d *Dissemination) HandleGossipMessage() {

}

func (d *Dissemination) SpreadMembershipListUpdates() {
}

func (d *Dissemination) send(nodeAddress string, gossipMessage []byte) {
	ctx, cancel := context.WithTimeout(g.gossipGlobalContext, g.timeoutTime)
	defer cancel()
	conn, err := net.Dial("tcp", nodeAddress)
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
		// Unmarshal Ack Messages
	}
}

func (d *Dissemination) marshalMembershipList(clusterData []*Node) ([]byte, error) {
	var (
		encodedMembershipList bytes.Buffer
		entry []byte
		err error
	)

	for index := range clusterData {
		entry, err = json.Marshal(map[string]any{
			"address": clusterData[index].nodeAddress,
			"port": clusterData[index].listenPort,
			"status": clusterData[index].nodeStatus,
		})

		if err != nil {
			return nil, err
		}

		encodedMembershipList.Write(entry)
	}

	return encodedMembershipList.Bytes(), nil
}