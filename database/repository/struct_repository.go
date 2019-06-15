package repository

// StructRepository is a repository interface for storing Struct types.
type StructRepository interface {
	Set(kvps ...KeyStructValue) Result
	// SetField updates a given struct's fields' values.
	SetField(entityType int, key string, fields map[string]interface{}) Result
	Get(entityType int, keys ...string) ([]KeyStructValue, Result)
	// GetField retrieves a given struct's fields' values.
	GetField(entityType int, key string, fields ...string) (KeyStructValue, Result)
	// Remove will delete the Structs from Store.
	Remove(entityType int, keys ...string) Result
}
