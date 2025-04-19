package vvector

import (
	id "github.com/google/uuid"
)

type VersionVector struct {
	dataVersion int64
}

func NewVersionVector() VersionVector {
	return VersionVector{
		dataVersion: 0,
	}
}

func (v *VersionVector) IncrementVector() {
	v.dataVersion += 1
}

func (v *VersionVector) UpdateVector(newVersion int64) {
	v.dataVersion = newVersion
}
