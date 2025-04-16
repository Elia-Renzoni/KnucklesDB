package vvector

import id "github.com/google/uuid"

type DataVersioning struct {
	version       int
	versionVector VersionVector
}

func NewDataVersioning(nodeId id.UUID) DataVersioning {
	return DataVersioning{
		version:       -1,
		versionVector: NewVersionVector(nodeId),
	}
}

func (d *DataVersioning) CompareAndUpdateVersions(receivedVersion VersionVector) {
	switch {
	case receivedVersion.dataVersion > d.versionVector.dataVersion:
		d.version = BEFORE
		d.versionVector.UpdateVector(receivedVersion.dataVersion)
	case receivedVersion.dataVersion < d.versionVector.dataVersion:
		d.version = AFTER
	case receivedVersion.dataVersion == d.versionVector.dataVersion:
		d.version = CONCURRENT
	}
}