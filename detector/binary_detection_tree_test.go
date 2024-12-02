package detector_test

import (
	"knucklesdb/detector"
	"testing"
)

func TestNewDetectionBST(t *testing.T) {
	if instance := detector.NewDetectionBST(); instance == nil {
		t.Fail()
	}
}

func TestInsertRoot(t *testing.T) {
	dTree := detector.NewDetectionBST()
	dTree.Insert("/foo", 78)

	if dTree.Root == nil {
		t.Fail()
	}
}

func TestInsert(t *testing.T) {
	dTree := detector.NewDetectionBST()
	dTree.Insert("/foo", 89)
	dTree.Insert("/bar", 112)
	dTree.Insert("/foobar", 55)

	node := dTree.Root
	for node != nil && node.GetNodeID() != "/bar" {
		if 112 < node.GetNodeLogicalClock() {
			node = node.GetNodeLeftChild()
		} else {
			node = node.GetNodeRightChild()
		}
	}

	if node == nil {
		t.Fail()
	}

	if node.GetNodeID() != "/bar" {
		t.Fail()
	}
}

func TestRemove(t *testing.T) {
	dTree := detector.NewDetectionBST()
	dTree.Insert("/oof", 50)
	dTree.Insert("/rab", 222)
	dTree.Insert("/oofrab", 21)

	dTree.Remove("/oofrab", 21)

	node := dTree.Root
	for node != nil && node.GetNodeID() != "/oofrab" {
		if 21 < node.GetNodeLogicalClock() {
			node = node.GetNodeLeftChild()
		} else {
			node = node.GetNodeRightChild()
		}
	}

	if node != nil {
		t.Fail()
	}
}