package repository

// KeyStructValue struct
type KeyStructValue struct {
	// (optional) Type is used to specify Entity type information or classification.
	// If needed and useful. You can default to 0, which may mean, unused or unclassified.
	Type int
	// Key of the Entity.
	Key   string
	// Value of the Entity. Value's struct data is serialized to a map type.
	Value map[string]interface{}
}

// NewKeyStructValue instantiates a Key StructValue object & set its members with received params.
func NewKeyStructValue(entityType int, key string, fields map[string]interface{}) *KeyStructValue {
	return &KeyStructValue{
		Type:  entityType,
		Key:   key,
		Value: fields,
	}
}
