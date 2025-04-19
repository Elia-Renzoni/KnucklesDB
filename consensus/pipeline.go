package consensus

import (
	"knucklesdb/vvector"
)

type PipelinedMessage struct {
	MessageType string `json:"type"`
	Pipeline []vvector.VersionVectorMessage `json:"data"`
}