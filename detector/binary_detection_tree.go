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

func (d *DetectionBST) Remove(nodeId string, clock int16) {
	var (
		node *TreeNode = d.root
		parent *TreeNode = d.root
		sub *TreeNode
	)

	for node != nil && node.value.nodeId != nodeId {
		parent = node
		if clock < node.value.clock {
			node = node.left
		} else {
			node = node.right
		}
	}

	if node != nil {
		if node.left == nil {
			if node == d.root {
				d.root = node.right
			} else {
				if clock < parent.value.logicalClock {
					parent.left = node.left
				} else {
					parent.right = node.right
				}
			}
		} else {
			if node.right == nil {
				if node == d.root {
					d.root = node.left
				} else {
					if clock < parent.value.logicalClock {
						parent.left = node.left
					} else {
						parent.right = node.left
					}
				}
			} else {
				sub = node
				parent = sub
				node = sub.left
				for node != nil {
					parent = node
					node = node.right
				}
				sub.value.nodeId = node.value.nodeId
				if parent == sub {
					parent.left = node.left
				} else {
					parent.right = node.left
				}
			}
		}
	}
}

func (d *DetectionBST) GetDLevels() int {
	return d.levels
}