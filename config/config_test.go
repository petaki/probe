package config

import (
	"os"
	"testing"
)

func TestLoadAndParse(t *testing.T) {
	os.Setenv("PROBE_REDIS_PREFIX", "probe:")
	os.Setenv("PROBE_REDIS_HOST", "127.0.0.1")
	os.Setenv("PROBE_REDIS_PASSWORD", "secret")
	os.Setenv("PROBE_REDIS_PORT", "6379")
	os.Setenv("PROBE_REDIS_DATABASE", "0")

	config, err := Load()

	if err != nil {
		t.Errorf("Can not load the environment variables.")
	}

	if config.RedisPrefix != "probe:" {
		t.Errorf("Expected redis prefix probe:, but got %v", config.RedisPrefix)
	}

	if config.RedisHost != "127.0.0.1" {
		t.Errorf("Expected redis host 127.0.0.1, but got %v", config.RedisHost)
	}

	if config.RedisPassword != "secret" {
		t.Errorf("Expected redis password secret, but got %v", config.RedisPassword)
	}

	if config.RedisPort != 6379 {
		t.Errorf("Expected redis port 6379, but got %v", config.RedisPort)
	}

	if config.RedisDatabase != 0 {
		t.Errorf("Expected redis database 0, but got %v", config.RedisDatabase)
	}
}
