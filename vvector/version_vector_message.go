package vvector

type VersionVectorMessage struct {
	Key []byte   `json:"key"`
	Value []byte `json:"value"`
	Version int `json:"version"`
}
