package detector

import (
	"knucklesdb/detector"
)

const estimateFaultPeriod int16 = 10

type FailureDetector struct {
	detectorTree *DetectionBST
	faultyNodes chan <-string
}

func NewFailureDetector(tree *DetectionBST) *FailureDetector {
	return &FailureDetector{
		detectorTree: tree,
		faultyNodes: make(chan <-string)
	}
}

func (f *FailureDetector) FaultDetection() {
	var (
		rootClock int16 = f.detectorTree.root.values.logicalClock
		sloppyClock int16
	)

	sloppyClock = rootClock - estimateFaultPeriod

	// binary search 
	for {
		if node := binarySearch(f.detectorTree.root, sloppyClock) node != nil {
			f.faultyNodes <- node.value.nodeId
			continue
		} else {
			break
		}
	} 

	close(f.faultyNodes)

	go removeFaultyNodes(f.detectorTree.root, f.faultyNodes)
}

func removeFaultyNodes(root *TreeNode, faultyNodes chan<-string) {
	for node := range faultyNodes {
		// TODO
	}
}

func binarySearch(root *TreeNode, key int16) *TreeNode {
	var node *TreeNode = root

	for node != nil && node.values.logicalClock > key {
		if key < node.values.logicalClock {
			node = node.left
		} else {
			node = node.right
		}
	}

	return node
}

