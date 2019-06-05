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
		if !repo.L1Cache.Upsert(kvps).IsSuccessful() {
			// delete from L1 cache so succeeding "gets" will reload from L2 Cache.
			repo.deleteFromL1Cache(kvps)
		}
		return e
	} else if e.Details != nil {
		failedUpserts := e.Details.([]UpsertFailDetail)
		if failedUpserts == nil || len(failedUpserts) == 0{
			return e
		}
		nkvps := make([]KeyValue,0, len(failedUpserts))
		for _,d := range kvps{
			// skip items that failed upsert as they are not persisted to L2 Cache
			if itemExists(d, failedUpserts){
				continue
			}
			nkvps = append(nkvps, d)
		}
		// sync L1 Cache with items that succeeded to L2 Cache upsert, intentionally ignore errors on L1 Cache.
		repo.L1Cache.Upsert(nkvps)
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
	} else if e.Details != nil {
		failedDeletes := e.Details.([]DeleteFailDetail)
		if failedDeletes == nil || len(failedDeletes) == 0{
			return e
		}
		keys := make([]string,0, len(failedDeletes))
		for _,k := range keys{
			// skip items that failed delete as they are not persisted to L2 Cache
			if itemKeyExists(k, failedDeletes){
				continue
			}
			keys = append(keys, k)
		}
		// sync L1 Cache with items that succeeded to L2 Cache delete, intentionally ignore errors on L1 Cache.
		repo.L1Cache.Delete(entityType, keys)
	}
	return e
}

func (repo l1l2Store)deleteFromL1Cache(kvps []KeyValue){
	keys := make([]string,0,len(kvps))
	var entityType int
	sameTypes := true
	for i,kvp := range kvps{
		keys = append(keys, kvp.Key)
		if i == 0 {
			entityType = kvp.Type
			continue
		} else if entityType != kvp.Type{
			sameTypes = false
		}
	}
	if sameTypes {
		repo.L1Cache.Delete(entityType, keys)
		return
	}
	for _,kvp := range kvps{
		repo.L1Cache.Delete(kvp.Type, []string{kvp.Key})
	}
	return
}

func itemExists(kvp KeyValue, kvps []UpsertFailDetail) bool{
	for i := range kvps{
		if kvp.Key == kvps[i].KeyValue.Key{
			return true
		}
	}
	return false
}
func itemKeyExists(key string, kvps []DeleteFailDetail) bool{
	for i := range kvps{
		if key == kvps[i].Key{
			return true
		}
	}
	return false
}
