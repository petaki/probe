package config

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestLoadAndParse(t *testing.T) {
	err := godotenv.Load("../.env.example")
	if err != nil {
		t.Errorf("Cannot load the .env.example file.")
	}

	config, err := Load()
	if err != nil {
		t.Errorf("Cannot load the environment variables.")
	}

	if config.RedisURL != "redis://127.0.0.1:6379/0" {
		t.Errorf("Expected redis host redis://127.0.0.1:6379/0, but got %v", config.RedisURL)
	}

	if config.RedisKeyPrefix != "probe:" {
		t.Errorf("Expected redis prefix probe:, but got %v", config.RedisKeyPrefix)
	}

	if config.RedisKeyTimeout != 2592000 {
		t.Errorf("Expected redis timeout 2592000, but got %v", config.RedisKeyTimeout)
	}
}
