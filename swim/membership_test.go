package swim_test

import (
	"testing"
	"knucklesdb/swim"
)

func TestIsSeed(t *testing.T) {
	clusterManager := swim.NewClusterManager()

	result, err := clusterManager.IsSeed("127.0.0.1", 5050)
	if err != nil {
		t.Log(err)
	}

	if result == false {
		t.Fail()
	}
}