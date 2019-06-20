package cache

// //import "time"
// import "github.com/SharedCode/parallels/database/repository"

// // MemoryCache is a repository that keeps data in memory,locally.
// type MemoryCache struct {
// 	lookup            map[string]repository.KeyValue
// 	durationInSeconds int
// }

// const defaultDuration = 30

// func NewMemoryCache(durationSecs int) repository.Repository {
// 	if durationSecs <= 0 {
// 		durationSecs = defaultDuration
// 	}
// 	return MemoryCache{
// 		lookup:            make(map[string]repository.KeyValue),
// 		durationInSeconds: durationSecs,
// 	}
// }
// func NewBigMemoryCache(size int, durationSecs int) repository.Repository {
// 	if durationSecs <= 0 {
// 		durationSecs = defaultDuration
// 	}
// 	return MemoryCache{
// 		lookup:            make(map[string]repository.KeyValue, size),
// 		durationInSeconds: durationSecs,
// 	}
// }

// // Set a set of KeyValue entries to the DB.
// func (repo MemoryCache) Set(kvps []repository.KeyValue) repository.Result {
// 	// for kvp := range kvps {
// 	// 	repo.lookup[kvp.]
// 	// }
// 	return repository.Result{}
// }

// // Get retrieves a set of KeyValue entries from DB given a set of Keys.
// func (repo MemoryCache) Get(entityType int, keys []string) ([]repository.KeyValue, repository.Result) {
// 	return nil, repository.Result{}
// }

// // Remove a set of entries in DB given a set of Keys.
// func (repo MemoryCache) Remove(entityType int, keys []string) repository.Result {
// 	return repository.Result{}
// }
