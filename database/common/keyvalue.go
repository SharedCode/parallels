package common

// KeyValue struct
type KeyValue struct {
	// (optional) Type is used to specify Entity type information or classification.
	// If needed and useful. Defaults to 0, which may mean, unused or unclassified.
	Type int
	// Key of the Entity.
	Key   string
	// Value of the Entity.
	Value []byte
}

// NewKeyValue instantiates a KeyValue object & set its members with received params.
func NewKeyValue(entityType int, key string, v []byte) *KeyValue {
	return &KeyValue{
		Type:  entityType,
		Key:   key,
		Value: v,
	}
}
