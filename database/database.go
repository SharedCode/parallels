package database

import "parallels/database/common"
import "parallels/database/cache"
import "parallels/database/store"

type cachedRepository struct {
	common.RepositorySet
	Cache common.Repository
}

// NewRepositorySet instantiates a new RepositorySet given a configuration.
func NewRepositorySet(configuration Configuration) (common.RepositorySet, error) {
	cr, e := newRepository(configuration)
	return cr.RepositorySet, e
}

// NewRepository instantiates a new Repository with Caching enabled.
func NewRepository(configuration Configuration) (common.Repository, error) {
	return newRepository(configuration)
}
func newRepository(configuration Configuration) (cachedRepository, error) {
	cc := cache.NewClient(configuration.RedisConfig)
	// The default Repository is "NOT" Navigable Repo! (false on Navigable)
	repo, e := store.NewRepository(configuration.CassandraConfig)
	if e != nil {
		return cachedRepository{}, e
	}
	repo2, e := store.NewNavigableRepository(configuration.CassandraConfig)
	cr := cachedRepository{
		Cache: cache.NewRepository(cc),
	}
	cr.Store = repo
	cr.NavigableStore = repo2
	cr.RepositorySet.Store = cr
	return cr, e
}

func (repo cachedRepository) Upsert(kvps []common.KeyValue) common.ResultStatus {
	// todo: implement "lock" based mgmt of Store and Cache later, for now,
	// we rely on Store mgmt result as the "synchronization"  fact.
	// i.e. - we rely on idea that if Store succeeds, THEN it is most likely, Cache mgmt
	// is safe to be done.
	e := repo.Store.Upsert(kvps)
	if e.IsSuccessful() {
		// update cache if Store was updated successfully.
		return repo.Cache.Upsert(kvps)
	}
	return e
}

func (repo cachedRepository) Get(entityType int, keys []string) ([]common.KeyValue, common.ResultStatus) {
	v, e := repo.Cache.Get(entityType, keys)
	if v != nil || !e.IsSuccessful() {
		return v, e
	}
	v, e = repo.Store.Get(entityType, keys)
	if v == nil && e.IsSuccessful() {
		return nil, common.ResultStatus{}
	}
	if e.IsSuccessful() {
		return v, repo.Cache.Upsert(v)
	}
	return v, e
}

func (repo cachedRepository) Delete(entityType int, keys []string) common.ResultStatus {
	// todo: implement "lock" based mgmt of Store and Cache later, for now,
	// we rely on Store mgmt result as the "synchronization"  fact.
	// i.e. - we rely on idea that if Store action succeeds, THEN most likely,
	// respective Cache mgmt action is safe to be done.
	e := repo.Store.Delete(entityType, keys)
	if e.IsSuccessful() {
		return repo.Cache.Delete(entityType, keys)
	}
	return e
}
