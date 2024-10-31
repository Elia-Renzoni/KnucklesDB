package detector_test

import (
	"knucklesdb/detector"
)

func TestNewDetectionBST(t *testing.T) {
	if instance := detector.NewDectionBST(); instance == nil {
		t.Fail()
	}
}