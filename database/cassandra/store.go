package cassandra

import "fmt"
import "time"
import "github.com/gocql/gocql"
import "github.com/SharedCode/parallels/database/repository"

type cassandraStore struct {
	connection Connection
	storeName  string
}

func NewNavigableRepository(config Config) (repository.NavigableRepository, error) {
	return newRepository(config, true)
}
func NewRepository(config Config) (repository.Repository, error) {
	return newRepository(config, false)
}

func CloseSession() {
	if globalSession == nil{return}
	globalSession.Close()
}

func (repo cassandraStore) Set(kvps ...repository.KeyValue) repository.Result {
	sql := fmt.Sprintf("UPDATE %s SET value=?, updated=?, is_del=false WHERE group=? AND key=?", repo.storeName)
	now := time.Now()

	session, e := repo.connection.getSession()
	if e != nil {
		return repository.Result{Error: e}
	}

	if repo.isStoreNavigable() {
		b := session.NewBatch(gocql.LoggedBatch)
		for _, kvp := range kvps {
			b.Query(sql, kvp.Value, now, kvp.Group, kvp.Key)
		}
		return repository.Result{Error: session.ExecuteBatch(b)}
	}
	// INSERT NOT using "batch" as batching in a "Key" that is a Partition Key, is anti-pattern(slows Cassandra down).
	var failedItems []repository.UpsertFailDetail
	for _, kvp := range kvps {
		kvp2 := kvp
		e := session.Query(sql, kvp2.Value, now, kvp2.Group, kvp2.Key).Exec()
		if e != nil {
			failedItems = append(failedItems, repository.UpsertFailDetail{KeyValue: kvp2, Error: e})
		}
	}

	if failedItems == nil || len(failedItems) == 0 {
		return repository.Result{}
	}
	return repository.Result{
		Error:        fmt.Errorf("Set failed upserting items, see ErrorDetails on which ones failed"),
		ErrorDetails: failedItems,
	}
}

func (repo cassandraStore) Get(group string, keys ...string) ([]repository.KeyValue, repository.Result) {

	session, e := repo.connection.getSession()
	if e != nil {
		return nil, repository.Result{Error: e}
	}

	inClause := ""
	for _, k := range keys {
		key := "'" + k + "'"
		if inClause == "" {
			inClause = key
			continue
		}
		inClause += ("," + key)
	}
	sql := fmt.Sprintf("SELECT key, value, is_del FROM %s WHERE group=? AND key IN ("+inClause+")", repo.storeName)
	iter := session.Query(sql, group).Iter()
	var kvps []repository.KeyValue
	m := map[string]interface{}{}
	for iter.MapScan(m) {
		if m["is_del"].(bool) {
			continue
		}
		if kvps == nil {
			kvps = make([]repository.KeyValue, 0, len(keys))
		}
		kvps = append(kvps, *repository.NewKeyValue(group, m["key"].(string), m["value"].([]byte)))
		m = map[string]interface{}{}
	}
	return kvps, repository.Result{}
}

func (repo cassandraStore) Remove(group string, keys ...string) repository.Result {
	sql := fmt.Sprintf("UPDATE %s SET updated=?, is_del=true WHERE group=? AND key=?", repo.storeName)
	now := time.Now()

	session, e := repo.connection.getSession()
	if e != nil {
		return repository.Result{Error: e}
	}

	if repo.isStoreNavigable() {
		b := session.NewBatch(gocql.LoggedBatch)
		for _, key := range keys {
			b.Query(sql, now, group, key)
		}
		return repository.Result{Error: session.ExecuteBatch(b)}
	}

	var failedItems []repository.DeleteFailDetail
	for _, key := range keys {
		e := session.Query(sql, now, group, key).Exec()
		if e != nil {
			failedItems = append(failedItems, repository.DeleteFailDetail{Key: key, Error: e})
		}
	}

	if failedItems == nil || len(failedItems) == 0{
		return repository.Result{}
	}
	return repository.Result{
		Error:        fmt.Errorf("Remove failed removing items, see ErrorDetails on which ones failed"),
		ErrorDetails: failedItems,
	}
}

func (repo cassandraStore) Navigate(group string, filter repository.Filter) ([]repository.KeyValue, repository.Result) {
	if !repo.isStoreNavigable() {
		return nil, repository.Result{Error: fmt.Errorf("Repository is not navigable")}
	}

	session, e := repo.connection.getSession()
	if e != nil {
		return nil, repository.Result{Error: e}
	}

	sql := "SELECT key, value, is_del FROM %s WHERE group=? AND key > ?"
	if filter.LessThanKey {
		sql = "SELECT key, value, is_del FROM %s WHERE group=? AND key < ?"
	}
	sql = fmt.Sprintf(sql, repo.storeName)
	iter := session.Query(sql, group, filter.Key).Iter()
	var kvps []repository.KeyValue
	m := map[string]interface{}{}
	for iter.MapScan(m) {
		if m["is_del"].(bool) {
			continue
		}
		kvps = append(kvps, repository.KeyValue{
			Group: group,
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
		connection: *c,
		storeName:  sn,
	}, e
}
