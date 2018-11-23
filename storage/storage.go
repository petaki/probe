package storage

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/petaki/probe/config"
)

// Storage type.
type Storage struct {
	Pool *redis.Pool
}

// New function.
func New(config config.Config) Storage {
	storage := Storage{}

	storage.Pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			options := []redis.DialOption{
				redis.DialDatabase(config.RedisDatabase),
			}

			if config.RedisPassword != "" {
				options = append(options, redis.DialPassword(config.RedisPassword))
			}

			return redis.Dial("tcp", config.RedisHost, options...)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}

			_, err := c.Do("PING")

			return err
		},
	}

	return storage
}
