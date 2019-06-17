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
	Repository
	// Navigate retrieves a set of data given its type classification (group)
	// and filter criteria. See Filter struct for filter definition.
	// Conditional filter such as Greater than or Less than Key is the typical use-case.
	Navigate(group string, filter Filter) ([]KeyValue, Result)
}
