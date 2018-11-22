package config

import "strconv"

// Current instance.
var Current = Config{}

// Config type.
type Config struct {
	RedisPrefix   string
	RedisHost     string
	RedisPassword string
	RedisPort     int
	RedisDatabase int
}

// GetRequiredKeys function.
func GetRequiredKeys() []string {
	return []string{
		"PROBE_REDIS_PREFIX",
		"PROBE_REDIS_HOST",
		"PROBE_REDIS_PASSWORD",
		"PROBE_REDIS_PORT",
		"PROBE_REDIS_DATABASE",
	}
}

// Parse function.
func (c *Config) Parse(key string, value string) error {
	if key == "PROBE_REDIS_PREFIX" {
		c.RedisPrefix = value
	} else if key == "PROBE_REDIS_HOST" {
		c.RedisHost = value
	} else if key == "PROBE_REDIS_PASSWORD" {
		c.RedisPassword = value
	} else if key == "PROBE_REDIS_PORT" {
		number, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}

		c.RedisPort = int(number)
	} else if key == "PROBE_REDIS_DATABASE" {
		number, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}

		c.RedisDatabase = int(number)
	}

	return nil
}
