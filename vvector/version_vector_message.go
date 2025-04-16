package vvector

type VersionVectorMessage struct {
	key []byte   `json:"key"`
	value []byte `json:"value"`
	version int `json:"version"`
}
