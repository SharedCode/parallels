package store

import "fmt"
import "time"
import "github.com/gocql/gocql"
import "github.com/SharedCode/parallels/database/repository"

type cassandraStore struct {
	Connection Connection
	storeName  string
}

func NewNavigableRepository(config Config) (repository.NavigableRepository, error) {
	return newRepository(config, true)
}
func NewRepository(config Config) (repository.Repository, error) {
	return newRepository(config, false)
}

func (repo cassandraStore) Set(kvps ...repository.GroupKeyValue) repository.Result {
	sql := fmt.Sprintf("UPDATE %s SET value=?, updated=?, is_del=false WHERE type=? AND key=?", repo.storeName)
	now := time.Now()

	if repo.isStoreNavigable() {
		b := repo.Connection.Session.NewBatch(gocql.LoggedBatch)
		for _, kvp := range kvps {
			b.Query(sql, kvp.Value, now, kvp.Type, kvp.Key)
		}
		return repository.Result{Error: repo.Connection.Session.ExecuteBatch(b)}
	}
	// INSERT NOT using "batch" as batching in a "Key" that is a Partition Key, is anti-pattern(slows Cassandra down).
	ch2 := make(chan error)
	// run Query in its own thread, to allow concurrent I/O to DB.
	for _, kvp := range kvps {
		kvp2 := kvp
		go func() { ch2 <- repo.Connection.Session.Query(sql, kvp2.Value, now, kvp2.Type, kvp2.Key).Exec() }()
	}
	var failedItems []repository.UpsertFailDetail
	// gather results
	for _, kvp := range kvps {
		e := <-ch2
		if e != nil {
			failedItems = append(failedItems, repository.UpsertFailDetail{GroupKeyValue: kvp, Error: e})
		}
	}
	if failedItems == nil {
		return repository.Result{}
	}
	return repository.Result{
		Error:   fmt.Errorf("Set failed upserting items, see Details on which ones failed"),
		Details: failedItems,
	}
}

func (repo cassandraStore) Get(entityType int, group string, keys ...string) ([]repository.GroupKeyValue, repository.Result) {
	inClause := ""
	for _, k := range keys {
		key := "'" + k + "'"
		if inClause == "" {
			inClause = key
			continue
		}
		inClause += ("," + key)
	}
	sql := fmt.Sprintf("SELECT key, value, is_del FROM %s WHERE type=? AND key IN ("+inClause+")", repo.storeName)
	iter := repo.Connection.Session.Query(sql, entityType).Iter()
	var kvps []repository.GroupKeyValue
	m := map[string]interface{}{}
	for iter.MapScan(m) {
		if m["is_del"].(bool) {
			continue
		}
		if kvps == nil {
			kvps = make([]repository.GroupKeyValue, 0, len(keys))
		}
		kvps = append(kvps, *repository.NewGroupKeyValue(entityType, group, m["key"].(string), m["value"].([]byte)))
		m = map[string]interface{}{}
	}
	return kvps, repository.Result{}
}

func (repo cassandraStore) Remove(entityType int, keys ...string) repository.Result {
	sql := fmt.Sprintf("UPDATE %s SET updated=?, is_del=true WHERE type=? AND key=?", repo.storeName)
	now := time.Now()
	if repo.isStoreNavigable() {
		b := repo.Connection.Session.NewBatch(gocql.LoggedBatch)
		for _, key := range keys {
			b.Query(sql, now, entityType, key)
		}
		return repository.Result{Error: repo.Connection.Session.ExecuteBatch(b)}
	}

	// scather exec query
	ch2 := make(chan error)
	for _, key := range keys {
		k2 := key
		go func() { ch2 <- repo.Connection.Session.Query(sql, now, entityType, k2).Exec() }()
	}

	// gather query results
	var failedItems []repository.DeleteFailDetail
	for _, key := range keys {
		e := <-ch2
		if e != nil {
			failedItems = append(failedItems, repository.DeleteFailDetail{Key: key, Error: e})
		}
	}

	if failedItems == nil {
		return repository.Result{}
	}
	return repository.Result{
		Error:   fmt.Errorf("Remove failed removing items, see Details on which ones failed"),
		Details: failedItems,
	}
}

func (repo cassandraStore) Navigate(entityType int, group string, filter repository.Filter) ([]repository.KeyValue, repository.Result) {
	if !repo.isStoreNavigable() {
		return nil, repository.Result{Error: fmt.Errorf("Repository is not navigable")}
	}
	sql := "SELECT key, value, is_del FROM %s WHERE type=? AND key > ?"
	if filter.LessThanKey {
		sql = "SELECT key, value, is_del FROM %s WHERE type=? AND key < ?"
	}
	sql = fmt.Sprintf(sql, repo.storeName)
	iter := repo.Connection.Session.Query(sql, entityType, filter.Key).Iter()
	var kvps []repository.KeyValue
	m := map[string]interface{}{}
	for iter.MapScan(m) {
		if m["is_del"].(bool) {
			continue
		}
		kvps = append(kvps, repository.KeyValue{
			Type:  entityType,
			Key:   m["key"].(string),
			Value: m["value"].([]byte),
		})
		m = map[string]interface{}{}
	}
	return kvps, repository.Result{}
}

var storeNameLiteral = "key_value"
var storeNameNavigableLiteral = "key_value_navigable"

func (repo cassandraStore) isStoreNavigable() bool {
	return repo.storeName == storeNameNavigableLiteral
}

func newRepository(config Config, navigableStore bool) (cassandraStore, error) {
	if config.TableName != "" {
		storeNameLiteral = config.TableName
	}
	if config.NavigableTableName != "" {
		storeNameNavigableLiteral = config.NavigableTableName
	}
	c, e := GetConnection(config)
	sn := storeNameLiteral
	if navigableStore {
		sn = storeNameNavigableLiteral
	}
	return cassandraStore{
		Connection: *c,
		storeName:  sn,
	}, e
}
