package cassandra;

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
	// NumCommns is Number of Connections per Host.
	NumConns int
	NavigableTableName string
	Port int
}

type Connection struct{
	Config
}

// GetConnection will create(& return) a new Connection to Cassandra if there is not one yet,
// otherwise, will just return existing singleton connection.
func GetConnection(config Config) (*Connection, error){
	if config.Keyspace == "" {
		return nil,fmt.Errorf("config.Keyspace is empty")
	}
	if config.NumConns <= 0{
		config.NumConns = 2
	}
	return &Connection{Config: config}, nil
}

var globalSession *gocql.Session
var locker sync.Mutex

func (conn *Connection)getSession() (*gocql.Session, error){
	if globalSession != nil && !globalSession.Closed(){
		return globalSession, nil
	}
	locker.Lock()
	defer locker.Unlock()
	if globalSession != nil && !globalSession.Closed(){
		return globalSession, nil
	}

	config := conn.Config
	cluster := gocql.NewCluster(config.ClusterHosts...)
	cluster.Keyspace = config.Keyspace
	pass := gocql.PasswordAuthenticator{Username: config.Username, Password: config.Password}
	cluster.Authenticator = pass
	cluster.NumConns = config.NumConns
	cluster.Consistency = gocql.ParseConsistency(config.Consistency)
	if config.Port > 0 {
		cluster.Port = config.Port
	}
	s,err := gocql.NewSession(*cluster)
	if err != nil {
		s.Close()
		return nil, err
	}
	globalSession = s

	return s,nil
}
