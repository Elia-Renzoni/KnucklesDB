package vvector

type DataVersioning struct {
	Order       int
}

func NewDataVersioning() *DataVersioning {
	return &DataVersioning{
		Order: -1,
	}
}

func (d *DataVersioning) CompareAndUpdateVersions(receivedVersion VersionVectorMessage, memorizedVersion int) {
	switch {
	case receivedVersion.Version > memorizedVersion:
		d.Order = HAPPENS_BEFORE
	case receivedVersion.Version < memorizedVersion:
		d.Order = HAPPENS_AFTER
	case receivedVersion.Version == memorizedVersion:
		d.Order = HAPPENS_CONCURRENT
	}
}