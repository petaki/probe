package config

import "testing"

func TestParse(t *testing.T) {
	config := Config{}

	config.parse("PROBE_REDIS_HOST", "127.0.0.1")
	config.parse("PROBE_REDIS_PORT", "6379")

	if config.RedisHost != "127.0.0.1" {
		t.Errorf("Expected redis host 127.0.0.1, but got %v", config.RedisHost)
	}

	if config.RedisPort != 6379 {
		t.Errorf("Expected redis port 6379, but got %v", config.RedisPort)
	}
}
