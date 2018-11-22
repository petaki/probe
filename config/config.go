package config

import "strconv"

// Current instance.
var Current = Config{}

// Config type.
type Config struct {
	redisPrefix   string
	redisHost     string
	redisPassword string
	redisPort     int
	redisDatabase int
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
		c.redisPrefix = value
	} else if key == "PROBE_REDIS_HOST" {
		c.redisHost = value
	} else if key == "PROBE_REDIS_PASSWORD" {
		c.redisPassword = value
	} else if key == "PROBE_REDIS_PORT" {
		number, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}

		c.redisPort = int(number)
	} else if key == "PROBE_REDIS_DATABASE" {
		number, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}

		c.redisDatabase = int(number)
	}

	return nil
}
