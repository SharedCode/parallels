package repository

// Filter specifies a given filter criteria, typically used for data retrieval/navigation.
type Filter struct {
	UpperboundKey string
	LowerboundKey string

	UpperboundKeyInclusive bool
	LowerboundKeyInclusive bool

	// MaxCountLimit specifies maximum number of matching records to be fetched from DB.
	// Default value (0) means there is no max, which, since will fetch all,
	// has to be used with caution not to overload with reading a lot of data from DB.
	MaxCountLimit int
}

// NavigableRepository is a Repository with add'l method allowing data retrieval
// via navigable expressions. See Filter struct for details on supported expression.
type NavigableRepository interface {
	Repository
	// Navigate retrieves a set of data given its data group
	// and filter criteria. See Filter struct for filter definition.
	Navigate(group string, filter Filter) ([]KeyValue, Result)
}
