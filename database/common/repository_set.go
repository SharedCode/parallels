package common

// Filter specifies a given filter criteria, typically used for data retrieval/navigation.
type Filter struct {
	// Key interested about.
	Key string
	// LessThanKey is a conditional filter, such as Greater than or Less than use-case.
	// Defaults to false to signify "greater than Key" filter condition.
	LessThanKey bool
}

// NavigableRepository is a Repository with add'l method allowing data retrieval via navigable expressions.
type NavigableRepository interface {
	Repository
	// Navigate retrieves a set of data given its type classification (entityType)
	// and filter criteria. See Filter struct for filter definition.
	// Conditional filter such as Greater than or Less than Key is the typical use-case.
	Navigate(entityType int, filter Filter) ([]KeyValue, ResultStatus)
}

// RepositorySet contains all available Repositories from a backend DB.
type RepositorySet struct {
	// Store is the default Entity Repository.
	Store Repository

	// NavigableStore is another Repository that entries of which, can be queried or filtered
	// for more advanced retrieval or rows' "navigation" use-case.
	NavigableStore NavigableRepository
}
