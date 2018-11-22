package config

import "testing"

func TestGetRequiredKeys(t *testing.T) {
	requiredKeys := GetRequiredKeys()

	if len(requiredKeys) != 5 {
		t.Errorf("Expected required keys length 5, but got %v", len(requiredKeys))
	}
}

func TestParse(t *testing.T) {
	config := Config{}

	config.Parse("PROBE_REDIS_HOST", "127.0.0.1")
	config.Parse("PROBE_REDIS_PORT", "6379")

	if config.redisHost != "127.0.0.1" {
		t.Errorf("Expected redis host 127.0.0.1, but got %v", config.redisHost)
	}

	if config.redisPort != 6379 {
		t.Errorf("Expected redis port 6379, but got %v", config.redisPort)
	}
}
