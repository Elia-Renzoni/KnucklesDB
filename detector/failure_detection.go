package detector

import (
	"knucklesdb/detector"
)

const estimateFaultPeriod int16 = 10

type FailureDetector struct {
	detectorTree *DetectionBST
}

func NewFailureDetector(tree *DetectionBST) *FailureDetector{
	return &FailureDetector{
		detectorTree: tree,
	}
}

func (f *FailureDetector) FaultDetection() {
	var (
		rootClock int16 = f.detectorTree.root.values.logicalClock
		sloppyClock int16
	)

	sloppyClock = rootClock - estimateFaultPeriod

	// binary search
}

