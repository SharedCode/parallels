package cache

import (
	"fmt"
	//"encoding/json"
	"time"
	"github.com/go-redis/redis"
)

// Options we send to(set a) Redis connection.
type Options struct{
	RedisOptions *redis.UniversalOptions
	DurationInSeconds int
}

func (opt Options) GetDuration() time.Duration{
	return time.Duration(opt.DurationInSeconds)*time.Second
}

type Connection struct{
	Client redis.UniversalClient
	Options Options
}

func DefaultOptions() Options {
	return Options{}
}

func newClient(options Options) Connection {
	client := redis.NewUniversalClient(options.RedisOptions)
	c := Connection{
		Client : client,
		Options: options,
	}
	return c
}

// Ping tests connectivity for redis (PONG should be returned)
func (connection Connection) Ping() error {
	pong, err := connection.Client.Ping().Result()
	if err != nil {
		return err
	}
	fmt.Println(pong, err)
	// Output: PONG <nil>

	return nil
}