package detector

type SingularUpdateQueue struct {
	updateQueue    chan *Victim
	detectorBuffer *DetectorBuffer
}

func NewSingularUpdateQueue(clockBuffer *DetectorBuffer) *SingularUpdateQueue {
	return &SingularUpdateQueue{
		updateQueue:    make(chan *Victim),
		detectorBuffer: clockBuffer,
	}
}

func (s *SingularUpdateQueue) AddVictimPage(page *Victim) {
	s.updateQueue <- page
}

func (s *SingularUpdateQueue) UpdateQueueReader() {
	for {
		select {
		case victim := <-s.updateQueue:
			ok := s.detectorBuffer.SearchPage(victim)
			if ok {
				s.detectorBuffer.UpdateVictimEpoch(victim)
			} else {
				s.detectorBuffer.AddToDetectorBuffer(victim)
			}
		}
	}
}
