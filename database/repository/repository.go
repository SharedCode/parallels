package repository

// Repository interface, a.k.a. - Data Store interface.
type Repository interface {
	// Set a set of KeyValue entries to the DB.
	Set(kvps ...KeyValue) Result
	// Get retrieves a set of KeyValue entries from DB
	// given a set of Keys.
	Get(group string, keys ...string) ([]KeyValue, Result)
	// Remove a set of entries in DB given a set of Keys.
	Remove(group string, keys ...string) Result
}