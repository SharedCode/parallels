package store

import "fmt"
import "time"
import "github.com/gocql/gocql"
import "parallels/database/common"

type cassandraStore struct {
	Connection Connection
	storeName  string
}

var storeNameLiteral = "key_value"
var storeNameNavigableLiteral = "key_value_navigable"

func NewNavigableRepository(config Config) (common.NavigableRepository, error) {
	return newRepository(config, true)
}
func NewRepository(config Config) (common.Repository, error) {
	return newRepository(config, false)
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

func (repo cassandraStore) isStoreNavigable() bool {
	return repo.storeName == storeNameNavigableLiteral
}

func (repo cassandraStore) Upsert(kvps []common.KeyValue) common.ResultStatus {
	sql := fmt.Sprintf("UPDATE %s SET value=?, updated=?, is_del=false WHERE type=? AND key=?", repo.storeName)
	now := time.Now()
	if repo.isStoreNavigable() {
		b := repo.Connection.Session.NewBatch(gocql.LoggedBatch)
		for _, kvp := range kvps {
			b.Query(sql, kvp.Value, now, kvp.Type, kvp.Key)
		}
		return common.ResultStatus{Error: repo.Connection.Session.ExecuteBatch(b)}
	}
	// INSERT NOT using "batch" as batching in a "Key" that is a Partition Key, is anti-pattern(slows Cassandra down).
	for _, kvp := range kvps {
		e := repo.Connection.Session.Query(sql, kvp.Value, now, kvp.Type, kvp.Key).Exec()
		if e != nil {
			// for now, it is all or nothing (if error, return error and fail rest of batch).
			return common.ResultStatus{Error: e}
		}
	}
	return common.ResultStatus{}
}

func (repo cassandraStore) Get(entityType int, keys []string) ([]common.KeyValue, common.ResultStatus) {
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
	var kvps []common.KeyValue
	m := map[string]interface{}{}
	for iter.MapScan(m) {
		if m["is_del"].(bool) {
			continue
		}
		if kvps == nil {
			kvps = make([]common.KeyValue, 0, len(keys))
		}
		kvps = append(kvps, common.KeyValue{
			Type:  entityType,
			Key:   m["key"].(string),
			Value: m["value"].([]byte),
		})
		m = map[string]interface{}{}
	}
	return kvps, common.ResultStatus{}
}

func (repo cassandraStore) Delete(entityType int, keys []string) common.ResultStatus {
	sql := fmt.Sprintf("UPDATE %s SET updated=?, is_del=true WHERE type=? AND key=?", repo.storeName)
	now := time.Now()
	if repo.isStoreNavigable() {
		b := repo.Connection.Session.NewBatch(gocql.LoggedBatch)
		for _, key := range keys {
			b.Query(sql, now, entityType, key)
		}
		return common.ResultStatus{Error: repo.Connection.Session.ExecuteBatch(b)}
	}
	for _, key := range keys {
		e := repo.Connection.Session.Query(sql, now, entityType, key).Exec()
		if e != nil {
			return common.ResultStatus{Error: e}
		}
	}
	return common.ResultStatus{}
}

func (repo cassandraStore) Navigate(entityType int, filter common.Filter) ([]common.KeyValue, common.ResultStatus) {
	if !repo.isStoreNavigable() {
		return nil, common.ResultStatus{Error: fmt.Errorf("Repository is not navigable.")}
	}

	sql := "SELECT key, value, is_del FROM %s WHERE type=? AND key > ?"
	if filter.LessThanKey {
		sql = "SELECT key, value, is_del FROM %s WHERE type=? AND key < ?"
	}

	sql = fmt.Sprintf(sql, repo.storeName)
	iter := repo.Connection.Session.Query(sql, entityType, filter.Key).Iter()
	var kvps []common.KeyValue
	m := map[string]interface{}{}
	for iter.MapScan(m) {
		if m["is_del"].(bool) {
			continue
		}
		kvps = append(kvps, common.KeyValue{
			Type:  entityType,
			Key:   m["key"].(string),
			Value: m["value"].([]byte),
		})
		m = map[string]interface{}{}
	}
	return kvps, common.ResultStatus{}
}
