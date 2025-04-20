package vvector

import id "github.com/google/uuid"

type DataVersioning struct {
	Order       int
}

func NewDataVersioning() *DataVersioning {
	return &DataVersioning{
		Order: -1,
	}
}

func (d *DataVersioning) CompareAndUpdateVersions(receivedVersion VersionVector, memorizedVersion int) {
	switch {
	case receivedVersion.dataVersion > memorizedVersion:
		d.Order = HAPPENS_BEFORE
	case receivedVersion.dataVersion < memorizedVersion:
		d.Order = HAPPENS_AFTER
	case receivedVersion.dataVersion == memorizedVersion:
		d.Order = HAPPENS_CONCURRENT
	}
}