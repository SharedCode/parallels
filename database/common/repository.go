package common

type ResultStatus struct{
	// Error is the summary level error.
	Error error
	// Details contains the details of the error. This is either transport or Storage system specific.
	Details interface{}
}

type UpsertFailDetail struct{
	KeyValue KeyValue
	Error error
}
type DeleteFailDetail struct{
	Key string
	Error error
}

// IsSuccessful returns false if result status signifies failure.
func (result ResultStatus) IsSuccessful() bool{
	return result.Error == nil
}

// Repository interface, a.k.a. - Data Store interface.
type Repository interface {
	// Upsert a set of KeyValue entries to the DB.
	Upsert(kvps []KeyValue) ResultStatus
	// Get retrieves a set of KeyValue entries from DB given a set of Keys.
	Get(entityType int, keys []string) ([]KeyValue, ResultStatus)
	// Delete a set of entries in DB given a set of Keys.
	Delete(entityType int, keys []string) ResultStatus
}
