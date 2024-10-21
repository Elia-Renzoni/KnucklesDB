package detector_test

import (
	"knucklesdb/detector"
)

func TestNode(t *testing.T) {
	node := NewTreeNode("", 44, nil, nil)
	if node == nil {
		t.Fail()
	}
}