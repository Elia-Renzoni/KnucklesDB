package vvector

import (
	"encoding/json"
)

type VersionVectorMarshaler struct {
}

func NewVersionVectorMarshaler() *VersionVectorMarshaler {
	return &VersionVectorMarshaler{}
}


func (v *VersionVectorMarshaler) MarshalVersionVectorMessage(key, value []byte, version int64) ([]byte, error) {
	var (
		encodedData []byte
		err error
	)

	encodedData, err = json.Marshal(map[string]any{
		"key": key,
		"value": value,
		"version": version,
	})

	return encodedData, err
}