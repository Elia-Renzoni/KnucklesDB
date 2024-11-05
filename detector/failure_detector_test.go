package detector_test


import (
	"knucklesdb/detector"
	"testing"
)

func TestNewFailureDetector(t *testing.T) {
	tree := detector.NewDectionBST()
	if instance := detector.NewFailureDetector(tree); instance == nil {
		t.Fail()
	}
}

func TestFaultDetection(t *testing.T) {
	dTree := detector.NewDectionBST()
	dTree.Insert("/foo", 50)
	dTree.Insert("/foo2", 45)
	dTree.Insert("/foo3", 40)
	dTree.Insert("/foo4", 39)
	dTree.Insert("/foo5", 20)

	faultsD := detector.NewFailureDetector(dTree)
	faultsD.FaultDetection()

	var removedItemsList []string = make([]string, 0)
	removedItemsList = append(removedItemsList, "/foo3")
	removedItemsList = append(removedItemsList, "/foo4")
	removedItemsList = append(removedItemsList, "/foo5")

	node := dTree.Root
	var (
		key int16
		removedCounter int = 0
	)
	for _, removedItem := range removedItemsList {
		switch removedItem {
		case "/foo3":
			key = 40
		case "/foo4":
			key = 39
		case "/foo5":
			key = 20
		}

		for node != nil && node.GetNodeID() != removedItem {
			if key < node.GetNodeLogicalClock() {
				node = node.GetNodeLeftChild()
			} else {
				node = node.GetNodeRightChild()
			}
		}
		if node == nil {
			removedCounter++
		}
	}

	t.Errorf("%d", removedCounter)
	
	if removedCounter == 3 {
		t.Fail()
	}
}
