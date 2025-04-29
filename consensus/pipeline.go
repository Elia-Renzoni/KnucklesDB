package consensus

import (
	"knucklesdb/vvector"

	id "github.com/google/uuid"
)

type PipelinedMessage struct {
	MessageType  string                         `json:"type"`
	ReplicaUUID  id.UUID                        `json:"uuid"`
	LogicalClock int                            `json:"clock"`
	Pipeline     []vvector.VersionVectorMessage `json:"data"`
}
