package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config type.
type Config struct {
	RedisPrefix   string
	RedisTimeout  int
	RedisHost     string
	RedisPassword string
	RedisPort     string
	RedisDatabase int
}

const (
	envRedisPrefix   = "PROBE_REDIS_PREFIX"
	envRedisTimeout  = "PROBE_REDIS_TIMEOUT"
	envRedisHost     = "PROBE_REDIS_HOST"
	envRedisPassword = "PROBE_REDIS_PASSWORD"
	envRedisPort     = "PROBE_REDIS_PORT"
	envRedisDatabase = "PROBE_REDIS_DATABASE"
)

var envKeys = []string{
	envRedisPrefix,
	envRedisTimeout,
	envRedisHost,
	envRedisPassword,
	envRedisPort,
	envRedisDatabase,
}

// Load function.
func Load() (Config, error) {
	config := Config{}

	for _, key := range envKeys {
		value, hasKey := os.LookupEnv(key)
		if !hasKey {
			return Config{}, fmt.Errorf("%v is not defined", key)
		}

		err := config.parse(key, value)
		if err != nil {
			return Config{}, err
		}
	}

	return config, nil
}

func (c *Config) parse(key string, value string) error {
	switch key {
	case envRedisPrefix:
		c.RedisPrefix = value
	case envRedisTimeout:
		number, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}

		c.RedisTimeout = int(number)
	case envRedisHost:
		c.RedisHost = value
	case envRedisPassword:
		c.RedisPassword = value
	case envRedisPort:
		c.RedisPort = value
	case envRedisDatabase:
		number, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}

		c.RedisDatabase = int(number)
	}

	return nil
}
