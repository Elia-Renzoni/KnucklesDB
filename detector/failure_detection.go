package detector

import (
	_"time"
	_"fmt"
)

const estimatedFaultPeriod int16 = 10

type FailureDetector struct {
	detectorTree *DetectionBST
	faultyNodes chan NodeValues
}

func NewFailureDetector(tree *DetectionBST) *FailureDetector {
	return &FailureDetector{
		detectorTree: tree,
		faultyNodes: make(chan NodeValues),
	}
}

func (f *FailureDetector) FaultDetection() {
	var (
		rootClock int16 = f.detectorTree.Root.GetNodeLogicalClock()
		sloppyClock int16
	)

	sloppyClock = rootClock - estimatedFaultPeriod

	go f.removeFaultyNodes()

	// binary search 
	go func () {
		for {
			if node := searchNode(f.detectorTree.Root, sloppyClock); node != nil {
				f.faultyNodes <- node.value
			} else {
				break
			}
		}	
		close(f.faultyNodes)
	}()
}

func (f *FailureDetector) removeFaultyNodes() {
	for {
		select {
		case node, ok := <- f.faultyNodes:
			if ok {
				f.detectorTree.Remove(node.nodeId, node.logicalClock)
			} else {
				break
			}
		}
	}
}

func searchNode(root *TreeNode, key int16) *TreeNode {
	var node *TreeNode = root

	for node != nil && node.GetNodeLogicalClock() > key {
		if key < node.GetNodeLogicalClock() {
			node = node.GetNodeLeftChild()
		} else {
			node = node.GetNodeRightChild()
		}
	}

	return node
}

