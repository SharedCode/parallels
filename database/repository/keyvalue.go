package repository

// KeyValue struct
type KeyValue struct {
	// (optional) Group is used to specify a data grouping construct.
	// Example, Entity Type, Time Series like hour, day, week value, etc...
	Group string
	// Key of the Entity.
	Key string
	// Value of the Entity.
	Value []byte
}

// NewKeyValue instantiates a KeyValue object & set its members with received params.
func NewKeyValue(group string, key string, v []byte) *KeyValue {
	return &KeyValue{
		Group: group,
		Key:   key,
		Value: v,
	}
}
