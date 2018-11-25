package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config type.
type Config struct {
	RedisHost       string
	RedisPassword   string
	RedisPort       string
	RedisDatabase   int
	RedisKeyPrefix  string
	RedisKeyTimeout int
}

const (
	envRedisHost       = "PROBE_REDIS_HOST"
	envRedisPassword   = "PROBE_REDIS_PASSWORD"
	envRedisPort       = "PROBE_REDIS_PORT"
	envRedisDatabase   = "PROBE_REDIS_DATABASE"
	envRedisKeyPrefix  = "PROBE_REDIS_KEY_PREFIX"
	envRedisKeyTimeout = "PROBE_REDIS_KEY_TIMEOUT"
)

var envKeys = []string{
	envRedisHost,
	envRedisPassword,
	envRedisPort,
	envRedisDatabase,
	envRedisKeyPrefix,
	envRedisKeyTimeout,
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
	case envRedisKeyPrefix:
		c.RedisKeyPrefix = value
	case envRedisKeyTimeout:
		number, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}

		c.RedisKeyTimeout = int(number)
	}

	return nil
}
