package detector

import (
	"time"
	"fmt"
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
		rootClock int16 = f.detectorTree.Root.value.logicalClock
		sloppyClock int16
	)

	sloppyClock = rootClock - estimatedFaultPeriod

	go f.removeFaultyNodes(f.detectorTree)

	// binary search 
	go func () {
		for {
			if node := searchNode(f.detectorTree.Root, sloppyClock); node != nil {
				f.faultyNodes <- node.value
				continue
			}
			break
		}	
	}()
 
	close(f.faultyNodes)
}

func (f *FailureDetector) removeFaultyNodes(tree *DetectionBST) {
	for {
		select {
		case node, ok := <- f.faultyNodes:
			if ok {
				tree.Remove(node.nodeId, node.logicalClock)
			} else {
				break
			}
		case <- time.After(5 *time.Second):
			fmt.Printf("timeout")
		default:
			fmt.Printf("...")
		}
	}
}

func searchNode(root *TreeNode, key int16) *TreeNode {
	var node *TreeNode = root

	for node != nil && node.value.logicalClock > key {
		if key < node.value.logicalClock {
			node = node.left
		} else {
			node = node.right
		}
	}

	fmt.Printf("%s", node.value.nodeId)

	return node
}

