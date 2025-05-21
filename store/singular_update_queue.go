/*	Copyright [2024] [Elia Renzoni]
*
*   Licensed under the Apache License, Version 2.0 (the "License");
*   you may not use this file except in compliance with the License.
*   You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*   Unless required by applicable law or agreed to in writing, software
*   distributed under the License is distributed on an "AS IS" BASIS,
*   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*   See the License for the specific language governing permissions and
*   limitations under the License.
*/



/**
* This file contains the implementation of the Singular Update Queue pattern,
*  according to which no locks are used to protect shared resources.
*  Instead, there is a single goroutine that reads from a channel and updates the information in the shared area itself.
*
*  channel [][][][][][][][][] <- reader -> buffer
 */

package store

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

/**
*  @brief This method must be called by the store to update the buffer with the new entries.
*  @param page to add
*
 */
func (s *SingularUpdateQueue) AddVictimPage(page *Victim) {
	s.updateQueue <- page
}

/**
*
* @brief This method implements the core of the goroutine that will read data from the buffer
*  and update the shared resource.
*
 */
func (s *SingularUpdateQueue) UpdateQueueReader() {
	for {
		// wait until the replacement algorithm is done
		s.detectorBuffer.wg.Wait()
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
