package repository

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

// GroupKeyValue is a KeyValue with Group added.
type GroupKeyValue struct{
	KeyValue
	// Group can be any string that conveys a grouping (construct) of data.
	// Samples are, time series (value representing hour, day, week) or some other
	// custom text that conveys a data group.
	// Groups are used in backend storage part of partition key (e.g. - in Cassandra).
	// This solves problems of "scale" I/O.
	Group string
}

func NewGroupKeyValue(entityType int, group string, key string, v []byte) *GroupKeyValue {
	r := GroupKeyValue{
		Group: group,
	}
	r.Type = entityType
	r.Key = key
	r.Value = v
	return &r
}
