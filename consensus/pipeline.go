package consensus

import (
	"knucklesdb/vvector"
)

type PipelinedMessage struct {
	Pipeline []vvector.VersionVectorMessage
}