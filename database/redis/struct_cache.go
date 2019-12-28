package redis

import "github.com/SharedCode/parallels/database/repository"

type structCache struct {
	redisConnection connection
}

// NewStructCache instantiates a cache Repository.
func NewStructCache(options Options) structCache {
	return structCache{
		redisConnection: newClient(options),
	}
}

// Set a set of entries to the cache.
func (repo structCache) Set(kvps ...repository.KeyStructValue) repository.Result {
	pipeline := repo.redisConnection.Client.Pipeline()
	for i := 0; i < len(kvps); i++ {
		for k,v := range kvps[i].Value{
			pipeline.HSet(format(kvps[i].Group, kvps[i].Key), k, v)
		}
	}
	// execute the batched upserts.
	cmdErr, e := pipeline.Exec()
	return extractError(e, cmdErr)
}

// // Get retrieves a set of entries from the cache.
// func (repo cache) Get(group string, keys ...string) ([]repository.KeyValue, repository.Result) {
// 	pipeline := repo.redisConnection.Client.Pipeline()
// 	m := map[string]*redis.StringCmd{}
// 	for i := 0; i < len(keys); i++ {
// 		m[keys[i]] = pipeline.Get(format(group, keys[i]))
// 	}
// 	// execute the batched upserts.
// 	cmdErr, e := pipeline.Exec()
// 	if e == redis.Nil {
// 		return nil, repository.Result{}
// 	}
// 	if e != nil {
// 		return nil, repository.Result{Error: e, ErrorDetails: cmdErr}
// 	}
// 	// process all returned results from the Server.
// 	var values []repository.KeyValue
// 	for k, v := range m {
// 		res, e := v.Result()
// 		if e != nil && e != redis.Nil {
// 			return nil, repository.Result{Error: e, ErrorDetails: cmdErr}
// 		}
// 		cmdErr = nil
// 		if values == nil {
// 			values = make([]repository.KeyValue, 0, len(keys))
// 		}
// 		values = append(values, *repository.NewKeyValue(group, k, []byte(res)))
// 	}
// 	return values, repository.Result{ErrorDetails: cmdErr}
// }

// // Remove a set of entries from the cache.
// func (repo cache) Remove(group string, keys ...string) repository.Result {
// 	pipeline := repo.redisConnection.Client.Pipeline()
// 	for i := 0; i < len(keys); i++ {
// 		pipeline.Del(format(group, keys[i]))
// 	}
// 	// execute the batched deletes.
// 	cmdErr, e := pipeline.Exec()
// 	return repository.Result{Error: e, ErrorDetails: cmdErr}
// }
