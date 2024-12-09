package detector

import (
	"fmt"
	"sync"
)

const estimatedFaultPeriod int16 = 10

type FailureDetector struct {
	detectorTree *DetectionBST
	faultyNodes chan NodeValues
	detectionWaiting chan struct{}
	helper *Helper
	wg sync.WaitGroup
}

func NewFailureDetector(tree *DetectionBST, helper *Helper, wg sync.WaitGroup, waiting chan struct{}) *FailureDetector {
	return &FailureDetector{
		detectorTree: tree,
		faultyNodes: make(chan NodeValues),
		detectionWaiting: waiting,
		helper: helper,
		wg: wg,
	}
}

func (f *FailureDetector) FaultDetection() {
	var (
		rootClock int16
		sloppyClock int16
	)
	for {	
		v := <- f.detectionWaiting
		fmt.Printf("%v", v)

		rootClock = f.detectorTree.Root.GetNodeLogicalClock()
		sloppyClock = rootClock - estimatedFaultPeriod

		go f.removeFaultyNodes()

		// binary search 
		go func (wg sync.WaitGroup) {
			defer wg.Done()
			wg.Add(1)

			for {
				if node := searchNode(f.detectorTree.Root, sloppyClock); node != nil {
					f.faultyNodes <- node.value
				} else {
					break
				}
			}	
			close(f.faultyNodes)
		}(f.wg)
	}
}

func (f *FailureDetector) removeFaultyNodes() {
	defer f.wg.Done()
	f.wg.Add(1)

	for {
		select {
		case node, ok := <- f.faultyNodes:
			if ok {
				f.detectorTree.Remove(node.nodeId, node.logicalClock)
				f.helper.AddNodeToEvict(node)
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

