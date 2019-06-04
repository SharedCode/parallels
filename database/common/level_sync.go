package common

type l1l2Store struct {
	L1Cache Repository
	L2Cache Repository
}

// NewL1L2Sync instantiates a new L1-L2 cache/store "synchronizer" as Repository.
func NewL1L2Sync(l1Cache Repository, l2Cache Repository) Repository {
	return l1l2Store{
		L1Cache: l1Cache,
		L2Cache: l2Cache,
	}
}

func (repo l1l2Store) Upsert(kvps []KeyValue) ResultStatus {
	// we rely on idea that if L2 Cache/Store succeeds, THEN it is most likely, mgmt
	// is safe to be done on L1 Cache. todo: prove this is NOT then implement "locking"!
	e := repo.L2Cache.Upsert(kvps)
	if e.IsSuccessful() {
		return repo.L1Cache.Upsert(kvps)
	}
	return e
}

func (repo l1l2Store) Get(entityType int, keys []string) ([]KeyValue, ResultStatus) {
	v, e := repo.L1Cache.Get(entityType, keys)
	if v != nil || !e.IsSuccessful() {
		return v, e
	}
	v, e = repo.L2Cache.Get(entityType, keys)
	if v == nil && e.IsSuccessful() {
		return nil, ResultStatus{}
	}
	if e.IsSuccessful() {
		return v, repo.L1Cache.Upsert(v)
	}
	return v, e
}

func (repo l1l2Store) Delete(entityType int, keys []string) ResultStatus {
	e := repo.L2Cache.Delete(entityType, keys)
	if e.IsSuccessful() {
		return repo.L1Cache.Delete(entityType, keys)
	}
	return e
}
