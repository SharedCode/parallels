package redis

import (
	"fmt"
	//"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

// Options we send to(set a) Redis connection.
type Options struct {
	RedisOptions      *redis.UniversalOptions
	DurationInSeconds int
}

func (opt Options) GetDuration() time.Duration {
	return time.Duration(opt.DurationInSeconds) * time.Second
}

// connection to Redis
type connection struct {
	Client  redis.UniversalClient
	Options Options
}

func newClient(options Options) connection {
	client := redis.NewUniversalClient(options.RedisOptions)
	c := connection{
		Client:  client,
		Options: options,
	}
	return c
}

// ping tests connectivity for redis (PONG should be returned)
func (connection connection) ping() error {
	pong, err := connection.Client.Ping().Result()
	if err != nil {
		return err
	}
	fmt.Println(pong, err)
	// Output: PONG <nil>

	return nil
}
