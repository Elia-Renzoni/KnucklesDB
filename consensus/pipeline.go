package consensus

import (
	"knucklesdb/vvector"

	id "github.com/google/uuid"
)

type PipelinedMessage struct {
	MessageType  string                         `json:"type"`
	RemoteAddr string `json:"remote_addr"`
	ReplicaUUID  id.UUID                        `json:"uuid"`
	LogicalClock int                            `json:"clock"`
	Pipeline     []vvector.VersionVectorMessage `json:"data"`
}
