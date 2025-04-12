package vvector

import (
	id "github.com/google/uuid"
)

type VersionVector struct {
	nodeID      id.UUID
	dataVersion int64
}

func NewVersionVector(node id.UUID) VersionVector {
	return VersionVector{
		nodeID:      node,
		dataVersion: 0,
	}
}

func (v *VersionVector) IncrementVector() {
	v.dataVersion += 1
}
