package repository

// StructRepository is a repository interface for storing Struct types.
type StructRepository interface {
	Set(kvps ...KeyStructValue) Result
	// SetField updates a given struct's fields' values.
	SetField(group string, key string, fields map[string]string) Result
	Get(group string, keys ...string) ([]KeyStructValue, Result)
	// GetField retrieves a given struct's fields' values.
	GetField(group string, key string, fields ...string) (KeyStructValue, Result)
	// Remove will delete the Structs from Store.
	Remove(group string, keys ...string) Result
}
