package detector

import (
	"knucklesdb/detector"
)

const estimatedFaultPeriod int16 = 10

type FailureDetector struct {
	detectorTree *DetectionBST
	faultyNodes chan <-NodeValues
}

func NewFailureDetector(tree *DetectionBST) *FailureDetector {
	return &FailureDetector{
		detectorTree: tree,
		faultyNodes: make(chan <-NodeValues)
	}
}

func (f *FailureDetector) FaultDetection() {
	var (
		rootClock int16 = f.detectorTree.root.values.logicalClock
		sloppyClock int16
	)

	sloppyClock = rootClock - estimatedFaultPeriod

	// binary search 
	for {
		if node := searchNode(f.detectorTree.root, sloppyClock) node != nil {
			f.faultyNodes <- node.value
			continue
		}
		break
	} 

	close(f.faultyNodes)

	go removeFaultyNodes(f.detectorTree, f.faultyNodes)
}

func removeFaultyNodes(tree *DetectionBST, faultyNodes chan<-NodeValues) {
	for node := range faultyNodes {
		tree.Remove(node.nodeId, node.logicalClock)
	}
}

func searchNode(root *TreeNode, key int16) *TreeNode {
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

