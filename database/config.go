package database

import (
	"github.com/go-redis/redis"
	rediscache "github.com/SharedCode/parallels/database/redis"
	"github.com/SharedCode/parallels/database/cassandra"
	"encoding/json"
	"io/ioutil")

// Configuration contains caching (redis) and backend store (e.g. Cassandra) host parameters.
type Configuration struct {
	RedisConfig     rediscache.Options
	CassandraConfig cassandra.Config
}

// LoadConfiguration will read from a JSON file the configuration & load it into memory.
func LoadConfiguration(filename string) (Configuration, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Configuration{}, err
	}

	var c Configuration
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return Configuration{}, err
	}
	// instantiates a Redis Universal Option that will connect to local host (default).
	if c.RedisConfig.RedisOptions == nil {
		c.RedisConfig.RedisOptions = &redis.UniversalOptions{}
	}
	return c, nil
}
