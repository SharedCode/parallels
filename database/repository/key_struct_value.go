package repository

// KeyStructValue struct
type KeyStructValue struct {
	// (optional) Group is used to specify Entity type information or classification.
	// If needed and useful. You can default to 0, which may mean, unused or unclassified.
	Group string
	// Key of the Entity.
	Key string
	// Value of the Entity. Value's struct data is serialized to a map type.
	Value map[string]string
}

// NewKeyStructValue instantiates a Key StructValue object & set its members with received params.
func NewKeyStructValue(group string, key string, fields map[string]string) *KeyStructValue {
	return &KeyStructValue{
		Group: group,
		Key:   key,
		Value: fields,
	}
}
