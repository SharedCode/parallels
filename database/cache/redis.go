package cache

import (
	"fmt"
	//"encoding/json"
	"time"
	"github.com/go-redis/redis"
)

// Options we send to(set a) Redis connection.
type Options struct{
	Address string
	Password string	
	DB int
	DurationInSeconds int
}

func (opt Options) GetDuration() time.Duration{
	return time.Duration(opt.DurationInSeconds)*time.Second
}

type Connection struct{
	Client *redis.Client
	Options Options
}

func DefaultOptions() Options {
	return Options{
		Address:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		DurationInSeconds: 0,	// no expiration!
	}
}

func NewClient(options Options) Connection {
	client := redis.NewClient(&redis.Options{
		Addr:     options.Address,
		Password: options.Password,
		DB:       options.DB})

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
