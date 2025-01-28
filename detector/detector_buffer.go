package detector

type DetectorBuffer struct {
	buffer     map[string]*Victim
	bufferSize int
}

func NewDetectorBuffer() *DetectorBuffer {
	return &DetectorBuffer{
		buffer: make(map[string]*Victim),
	}
}

func (d *DetectorBuffer) AddToDetectorBuffer(v *Victim) {
	d.buffer[string(v.key)] = v
	d.bufferSize = len(d.buffer)
}

func (d *DetectorBuffer) SearchPage(v *Victim) bool {
	_, ok := d.buffer[string(v.key)]
	if ok {
		return ok
	}
	return false
}

func (d *DetectorBuffer) UpdateVictimEpoch(v *Victim) {
	d.buffer[string(v.key)] = v
}
