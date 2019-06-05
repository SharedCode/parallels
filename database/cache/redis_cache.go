package cache

import "fmt"
import "github.com/go-redis/redis"
import "parallels/database/common"

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

func format(entityType int, key string) string { return fmt.Sprintf("%d_%s", entityType, key) }

// Upsert a set of entries to the cache.
func (repo RedisCache) Upsert(kvps []common.KeyValue) common.Result {
	pipeline := repo.redisConnection.Client.Pipeline()
	expiration := repo.redisConnection.Options.GetDuration()
	for i := 0; i < len(kvps); i++ {
		pipeline.Set(format(kvps[i].Type, kvps[i].Key), kvps[i].Value, expiration)
	}
	// execute the batched upserts.
	cmdErr, e := pipeline.Exec()
	return extractError(e, cmdErr)
}

func extractError(e error, cmdErr []redis.Cmder) common.Result {
	if e == nil && cmdErr != nil && len(cmdErr) >= 1 {
		if len(cmdErr) == 1 {
			err := cmdErr[0].Err()
			if err == nil {
				return common.Result{}
			}
		}
		e = fmt.Errorf("Error was encountered while working on the batch. See Details for more info")
	}
	return common.Result{Error: e, Details: cmdErr}
}

// Get retrieves a set of entries from the cache.
func (repo RedisCache) Get(entityType int, keys []string) ([]common.KeyValue, common.Result) {
	pipeline := repo.redisConnection.Client.Pipeline()
	m := map[string]*redis.StringCmd{}
	for i := 0; i < len(keys); i++ {
		m[keys[i]] = pipeline.Get(format(entityType, keys[i]))
	}
	// execute the batched upserts.
	cmdErr, e := pipeline.Exec()
	if e == redis.Nil {
		return nil, common.Result{}
	}
	if e != nil {
		return nil, common.Result{Error: e, Details: cmdErr}
	}
	// process all returned results from the Server.
	var values []common.KeyValue
	for k, v := range m {
		res, e := v.Result()
		if e != nil && e != redis.Nil {
			return nil, common.Result{Error: e, Details: cmdErr}
		}
		cmdErr = nil
		if values == nil {
			values = make([]common.KeyValue, 0, len(keys))
		}
		values = append(values, *common.NewKeyValue(entityType, k, []byte(res)))
	}
	return values, common.Result{Details: cmdErr}
}

// Delete a set of entries from the cache.
func (repo RedisCache) Delete(entityType int, keys []string) common.Result {
	pipeline := repo.redisConnection.Client.Pipeline()
	for i := 0; i < len(keys); i++ {
		pipeline.Del(format(entityType, keys[i]))
	}
	// execute the batched deletes.
	cmdErr, e := pipeline.Exec()
	return common.Result{Error: e, Details: cmdErr}
}
