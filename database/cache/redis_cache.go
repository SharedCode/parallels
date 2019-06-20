package cache

import "fmt"
import "github.com/go-redis/redis"
import "github.com/SharedCode/parallels/database/repository"

// RedisCache implementation for Redis based caching.
type RedisCache struct {
	redisConnection Connection
}

// NewRedisCache instantiates a cache Repository.
func NewRedisCache(options Options) RedisCache {
	return RedisCache{
		redisConnection: newClient(options),
	}
}

// Set a set of entries to the cache.
func (repo RedisCache) Set(kvps ...repository.KeyValue) repository.Result {
	pipeline := repo.redisConnection.Client.Pipeline()
	expiration := repo.redisConnection.Options.GetDuration()
	for i := 0; i < len(kvps); i++ {
		pipeline.Set(format(kvps[i].Group, kvps[i].Key), kvps[i].Value, expiration)
	}
	// execute the batched upserts.
	cmdErr, e := pipeline.Exec()
	return extractError(e, cmdErr)
}

// Get retrieves a set of entries from the cache.
func (repo RedisCache) Get(group string, keys ...string) ([]repository.KeyValue, repository.Result) {
	pipeline := repo.redisConnection.Client.Pipeline()
	m := map[string]*redis.StringCmd{}
	for i := 0; i < len(keys); i++ {
		m[keys[i]] = pipeline.Get(format(group, keys[i]))
	}
	// execute the batched upserts.
	cmdErr, e := pipeline.Exec()
	if e == redis.Nil {
		return nil, repository.Result{}
	}
	if e != nil {
		return nil, repository.Result{Error: e, ErrorDetails: cmdErr}
	}
	// process all returned results from the Server.
	var values []repository.KeyValue
	for k, v := range m {
		res, e := v.Result()
		if e != nil && e != redis.Nil {
			return nil, repository.Result{Error: e, ErrorDetails: cmdErr}
		}
		cmdErr = nil
		if values == nil {
			values = make([]repository.KeyValue, 0, len(keys))
		}
		values = append(values, *repository.NewKeyValue(group, k, []byte(res)))
	}
	return values, repository.Result{ErrorDetails: cmdErr}
}

// Remove a set of entries from the cache.
func (repo RedisCache) Remove(group string, keys ...string) repository.Result {
	pipeline := repo.redisConnection.Client.Pipeline()
	for i := 0; i < len(keys); i++ {
		pipeline.Del(format(group, keys[i]))
	}
	// execute the batched deletes.
	cmdErr, e := pipeline.Exec()
	return repository.Result{Error: e, ErrorDetails: cmdErr}
}

func extractError(e error, cmdErr []redis.Cmder) repository.Result {
	if e == nil && cmdErr != nil && len(cmdErr) >= 1 {
		if len(cmdErr) == 1 {
			err := cmdErr[0].Err()
			if err == nil {
				return repository.Result{}
			}
		}
		e = fmt.Errorf("Error was encountered while working on the batch. See ErrorDetails for more info")
	}
	return repository.Result{Error: e, ErrorDetails: cmdErr}
}

func format(group string, key string) string { return fmt.Sprintf("%d_%s", group, key) }
