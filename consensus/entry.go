package consensus

type Entry struct {
	key, value []byte
	version int64
}

func NewEntry(key, value []byte, version int64) Entry {
	return Entry{
		key: key,
		value: value,
		version: version,
	}
}