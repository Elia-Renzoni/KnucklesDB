package vvector

import id "github.com/google/uuid"

type DataVersioning struct {
	order       int
	versionVector VersionVector
}

func NewDataVersioning() *DataVersioning {
	return &DataVersioning{
		order:       -1,
		versionVector: NewVersionVector(nodeId),
	}
}

func (d *DataVersioning) CompareAndUpdateVersions(receivedVersion VersionVector) {
	switch {
	case receivedVersion.dataVersion > d.versionVector.dataVersion:
		d.order = BEFORE
		d.versionVector.UpdateVector(receivedVersion.dataVersion)
	case receivedVersion.dataVersion < d.versionVector.dataVersion:
		d.order = AFTER
	case receivedVersion.dataVersion == d.versionVector.dataVersion:
		d.order = CONCURRENT
	}
}