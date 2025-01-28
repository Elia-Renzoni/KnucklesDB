/**
*	This file contains the structure of a cache page stored in the buffer queue
*	for the clock algorithm
*
**/

package detector

type Victim struct {
	epoch  bool
	key    []byte
	pageID int
}

func NewVictim(key []byte, pid int) *Victim {
	return &Victim{
		epoch:  true,
		key:    key,
		pageID: pid,
	}
}
