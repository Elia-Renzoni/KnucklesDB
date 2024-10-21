package detector

type TreeNode struct {
	value NodeValues
	left *TreeNode
	right *TreeNode
}

type NodeValues struct {
	nodeId string
	logicalClock int16
}

func NewTreeNode(id string, clock int16, left, right *TreeNode) *TreeNode {
	return &TreeNode{
		value: NodeValues{
			nodeId: id,
			logicalClock: clock,
		},
		left: left,
		righ: right,
	}
}