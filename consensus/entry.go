package consensus

type Entry struct {
	key, value []byte
	version int
}

func NewEntry(key, value []byte, version int) Entry {
	return Entry{
		key: key,
		value: value,
		version: version,
	}
}