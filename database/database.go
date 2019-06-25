package database

import (
	"github.com/SharedCode/parallels/database/redis"
	"github.com/SharedCode/parallels/database/cassandra"
	"github.com/SharedCode/parallels/database/repository"
)

// NewRepositorySet instantiates a new RepositorySet given a configuration.
func NewRepositorySet(configuration Configuration) (repository.RepositorySet, error) {
	s1, e := NewRepository(configuration)
	if e != nil {
		return repository.RepositorySet{}, e
	}
	s2, e := cassandra.NewNavigableRepository(configuration.CassandraConfig)
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
	repo, e := cassandra.NewRepository(configuration.CassandraConfig)
	if e != nil {
		return nil, e
	}
	cache := redis.NewCache(configuration.RedisConfig)
	return repository.NewL1L2CacheSync(cache, repo), nil
}

// CloseSession closes all global Sessions that are opened.
func CloseSession(){
	cassandra.CloseSession()
}
