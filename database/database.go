package database

import "parallels/database/common"
import "parallels/database/cache"
import "parallels/database/store"

// NewRepositorySet instantiates a new RepositorySet given a configuration.
func NewRepositorySet(configuration Configuration) (common.RepositorySet, error) {
	l1,e := NewRepository(configuration)
	if e != nil{
		return common.RepositorySet{},e
	}
	l2,e := store.NewNavigableRepository(configuration.CassandraConfig)
	if e != nil{
		return common.RepositorySet{},e
	}
	return common.RepositorySet{
		Store: l1,
		NavigableStore: l2,
	}, nil
}

// NewRepository instantiates a new Repository with Caching enabled.
func NewRepository(configuration Configuration) (common.Repository, error) {
	repo, e := store.NewRepository(configuration.CassandraConfig)
	if e != nil {
		return nil, e
	}
	cache := cache.NewRedisCache(configuration.RedisConfig)
	return common.NewL1L2Sync(cache, repo), nil
}
