package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config type.
type Config struct {
	RedisUrl        string
	RedisKeyPrefix  string
	RedisKeyTimeout int
}

const (
	envRedisUrl        = "PROBE_REDIS_URL"
	envRedisKeyPrefix  = "PROBE_REDIS_KEY_PREFIX"
	envRedisKeyTimeout = "PROBE_REDIS_KEY_TIMEOUT"
)

var envKeys = []string{
	envRedisUrl,
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
	case envRedisUrl:
		c.RedisUrl = value
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
