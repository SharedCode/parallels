package repository

// Filter specifies a given filter criteria, typically used for data retrieval/navigation.
type Filter struct {
	// Key interested about.
	Key string
	// LessThanKey is a conditional filter, such as Greater than or Less than use-case.
	// Defaults to false to signify "greater than Key" filter condition.
	LessThanKey bool
}

// NavigableRepository is a Repository with add'l method allowing data retrieval 
// via navigable expressions.
type NavigableRepository interface {
	// Set a set of KeyValue entries to the DB.
	Set(kvps ...GroupKeyValue) Result
	// Get retrieves a set of KeyValue entries from DB
	// given a set of Keys.
	Get(entityType int, group string, keys ...string) ([]GroupKeyValue, Result)
	// Remove a set of entries in DB given a set of Keys.
	Remove(entityType int, group string, keys ...string) Result
	// Navigate retrieves a set of data given its type classification (entityType)
	// and filter criteria. See Filter struct for filter definition.
	// Conditional filter such as Greater than or Less than Key is the typical use-case.
	Navigate(entityType int, group string, filter Filter) ([]GroupKeyValue,Result)
}
