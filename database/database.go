package database

import "github.com/SharedCode/parallels/database/repository"
import "github.com/SharedCode/parallels/database/cache"
import "github.com/SharedCode/parallels/database/store"

// NewRepositorySet instantiates a new RepositorySet given a configuration.
func NewRepositorySet(configuration Configuration) (repository.RepositorySet, error) {
	s1, e := NewRepository(configuration)
	if e != nil {
		return repository.RepositorySet{}, e
	}
	s2, e := store.NewNavigableRepository(configuration.CassandraConfig)
	if e != nil {
		return repository.RepositorySet{}, e
	}
	return repository.RepositorySet{
		Store:          s1,
		NavigableStore: s2,
	}, nil
}

// NewRepository instantiates a new Repository with Caching enabled.
func NewRepository(configuration Configuration) (repository.Repository, error) {
	repo, e := store.NewRepository(configuration.CassandraConfig)
	if e != nil {
		return nil, e
	}
	cache := cache.NewRedisCache(configuration.RedisConfig)
	return repository.NewL1L2CacheSync(cache, repo), nil
}
