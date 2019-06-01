
package store;

import "fmt"
import "sync"
import "github.com/gocql/gocql"

type Config struct{
	ClusterHosts []string
	// Keyspace to be used when doing I/O to cassandra.
	Keyspace string
	Username string
	Password string
	TableName string
	Consistency string
	NavigableTableName string
	Port int
}

type Connection struct{
	Session *gocql.Session
	Config
}

var connection *Connection
var mux sync.Mutex

// GetConnection will create(& return) a new Connection to Cassandra if there is not one yet,
// otherwise, will just return existing singleton connection.
func GetConnection(config Config) (*Connection, error){
	if connection != nil && connection.Session != nil && !connection.Session.Closed(){
		return connection, nil
	}
	mux.Lock()
	defer mux.Unlock()

	if connection != nil {
		return connection, nil
	}
	if config.Keyspace == "" {
		return nil,fmt.Errorf("config.Keyspace is empty")
	}
	cluster := gocql.NewCluster(config.ClusterHosts...)
	cluster.Keyspace = config.Keyspace
	pass := gocql.PasswordAuthenticator{Username: config.Username, Password: config.Password}
	cluster.Authenticator = pass
	cluster.NumConns = 1
	cluster.Consistency = gocql.ParseConsistency(config.Consistency)
	if config.Port > 0 {
		cluster.Port = config.Port
	}
	var c = Connection{
		Config: config,
	}
	s, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	c.Session = s
	connection = &c
	return connection, nil
}
