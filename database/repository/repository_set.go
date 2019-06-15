package repository

// RepositorySet contains all available Repositories from a backend DB.
type RepositorySet struct {
	// Store is the default Entity Repository.
	Store Repository

	// NavigableStore is another Repository that entries of which, can be queried or filtered
	// for more advanced retrieval or rows' "navigation" use-case.
	NavigableStore NavigableRepository
}
