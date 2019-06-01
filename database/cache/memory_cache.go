package cache

import "parallels/database/common"

// Repository interface, a.k.a. - Data Store interface.
type MemoryCache struct {

}

func NewMemoryCache() common.Repository{
	return MemoryCache{}
}

// Upsert a set of KeyValue entries to the DB.
func(repo MemoryCache) Upsert(kvps []common.KeyValue) common.ResultStatus{
	return common.ResultStatus{}
}
// Get retrieves a set of KeyValue entries from DB given a set of Keys.
func(repo MemoryCache) Get(entityType int, keys []string) ([]common.KeyValue, common.ResultStatus){
	return nil, common.ResultStatus{}
}
// Delete a set of entries in DB given a set of Keys.
func(repo MemoryCache) Delete(entityType int, keys []string) common.ResultStatus{
	return common.ResultStatus{}
}
