package swim_test

import (
	"testing"
	"knucklesdb/swim"
)


func TestMarshalPing(t *testing.T) {
	marshaler := swim.NewProtocolMarshaler()
	jsonValue, err := marshaler.MarshalPing()

	if err != nil {
		t.Fail()
	}

	t.Log(string(jsonValue))
} 

func TestMarshalPiggyBack(t *testing.T) {
	marshaler := swim.NewProtocolMarshaler()
	jsonValue, err := marshaler.MarshalPiggyBack("192.89.33.4", "192.244.66.77")

	if err != nil {
		t.Fail()
	}

	t.Log(string(jsonValue))
}

func TestMarshalSWIMDetectionMessage(t *testing.T) {
	marshaler := swim.NewProtocolMarshaler()
	jsonValue, err := marshaler.MarshalSWIMDetectionMessage(1, 4040, "127.0.0.1")

	if err != nil {
		t.Fail()
	}

	t.Log(string(jsonValue))
}