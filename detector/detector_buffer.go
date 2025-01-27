package detector


type DetectorBuffer struct {
	buffer []Victim	
}

func NewDetectorBuffer() *DetectorBuffer {
	return &DetectorBuffer{
		buffer: make([]Victim, 0)
	}
}

func (d *DetectorBuffer) AddToDetectorBuffer(v Victim) {

}