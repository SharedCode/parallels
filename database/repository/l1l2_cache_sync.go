package repository

import "fmt"

type l1l2Store struct {
	L1Cache Repository
	L2Cache Repository
}

// NewL1L2CacheSync instantiates a new L1-L2 cache/store "synchronizer" as Repository.
func NewL1L2CacheSync(l1Cache Repository, l2Cache Repository) Repository {
	return l1l2Store{
		L1Cache: l1Cache,
		L2Cache: l2Cache,
	}
}

func (repo l1l2Store) Set(kvps ...KeyValue) Result {
	// we rely on idea that if L2 Cache/Store succeeds, THEN it is most likely, mgmt
	// is safe to be done on L1 Cache. todo: prove this is NOT then implement "locking"!
	e := repo.L2Cache.Set(kvps...)
	if e.IsSuccessful() {
		e2 := repo.L1Cache.Set(kvps...)
		if !e2.IsSuccessful() {
			// delete from L1 cache so succeeding "gets" will reload from L2 Cache.
			// intentionally ignore errors on L1 cache delete, for now.
			repo.deleteFromL1Cache(kvps...)
		}
		return e
	} else if e.Details != nil {
		failedUpserts := e.Details.([]UpsertFailDetail)
		if failedUpserts == nil || len(failedUpserts) == 0 {
			return e
		}
		nkvps := make([]KeyValue, 0, len(failedUpserts))
		for _, d := range kvps {
			// skip items that failed upsert as they are not persisted to L2 Cache
			if itemExists(d, failedUpserts) {
				continue
			}
			nkvps = append(nkvps, d)
		}
		// sync L1 Cache with items that succeeded to L2 Cache upsert,
		// intentionally ignore errors on L1 Cache.
		repo.L1Cache.Set(nkvps...)
	}
	return e
}

func (repo l1l2Store) Get(entityType int, keys ...string) ([]KeyValue, Result) {
	kvps, result := repo.L1Cache.Get(entityType, keys...)
	if kvps != nil || !result.IsSuccessful() {
		return kvps, result
	}
	kvps, result = repo.L2Cache.Get(entityType, keys...)
	if kvps == nil && result.IsSuccessful() {
		return nil, Result{}
	}
	if result.IsSuccessful() {
		// sync up L1 cache.
		// todo: do we want to handle error on L1 cache ? prove it, then prolly remove from cache the "set"..
		repo.L1Cache.Set(kvps...)
	}
	return kvps, result
}

func (repo l1l2Store) Remove(entityType int, keys ...string) Result {
	result := repo.L2Cache.Remove(entityType, keys...)
	if result.IsSuccessful() {
		repo.L1Cache.Remove(entityType, keys...)
	} else if result.Details != nil {
		failedDeletes := result.Details.([]DeleteFailDetail)
		if failedDeletes == nil || len(failedDeletes) == 0 {
			return result
		}
		nkeys := make([]string, 0, len(failedDeletes))
		for _, k := range keys {
			// skip items that failed delete as they are not persisted to L2 Cache
			if itemKeyExists(k, failedDeletes) {
				continue
			}
			nkeys = append(nkeys, k)
		}
		// sync L1 Cache with items that succeeded to L2 Cache delete,
		// intentionally ignore errors on L1 Cache.
		repo.L1Cache.Remove(entityType, nkeys...)
	}
	return result
}

func (repo l1l2Store) deleteFromL1Cache(kvps ...KeyValue) Result {
	keys := make([]string, 0, len(kvps))
	var entityType int
	sameTypes := true
	for i, kvp := range kvps {
		keys = append(keys, kvp.Key)
		if i == 0 {
			entityType = kvp.Type
			continue
		} else if entityType != kvp.Type {
			sameTypes = false
		}
	}
	if sameTypes {
		return repo.L1Cache.Remove(entityType, keys...)
	}
	errors := make([]Result, len(kvps))
	for _, kvp := range kvps {
		r := repo.L1Cache.Remove(kvp.Type, kvp.Key)
		if !r.IsSuccessful() {
			errors = append(errors, r)
		}
	}
	if len(errors) == 0 {
		return Result{}
	}
	return Result{Error: fmt.Errorf("Remove from cache encountered failure, see Result.Details"), Details: errors}
}

func itemExists(kvp KeyValue, kvps []UpsertFailDetail) bool {
	for i := range kvps {
		if kvp.Key == kvps[i].KeyValue.Key {
			return true
		}
	}
	return false
}
func itemKeyExists(key string, kvps []DeleteFailDetail) bool {
	for i := range kvps {
		if key == kvps[i].Key {
			return true
		}
	}
	return false
}
