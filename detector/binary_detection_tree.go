package detector

import (
	"knucklesdb/detector"
)

type DetectionBST struct {
	root *TreeNode
	levels int
}

func NewDectionBST() *DetectionBST {
	return &DetectionBST{
		root: nil,
		levels: 0,
	}
}

func (d *DetectionBST) Insert(nodeId string, clock int16) {
	var (
		node *TreeNode = d.root
		parent *TreeNode = d.root
		newNode *TreeNode
	)

	for node != nil && node.value.clock != clock {
		parent = node
		if clock < node.value.clock {
			node = node.left
		} else {
			node = node.right
		}
	}

	if node == nil {
		newNode = detector.NewTreeNode(nodeId, clock, nil, nil)
		if node == d.root {
			d.root = newNode
		} else {
			if clock < parent.value.logicalClock {
				parent.left = newNode
			} else {
				parent.right = newNode
			}
		}
	}
}

func (d *DetectionBST) GetDLevels() int {
	return d.levels
}