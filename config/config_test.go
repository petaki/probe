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

	if config.RedisHost != "127.0.0.1" {
		t.Errorf("Expected redis host 127.0.0.1, but got %v", config.RedisHost)
	}

	if config.RedisPassword != "" {
		t.Errorf("Expected redis password empty, but got %v", config.RedisPassword)
	}

	if config.RedisPort != "6379" {
		t.Errorf("Expected redis port 6379, but got %v", config.RedisPort)
	}

	if config.RedisDatabase != 0 {
		t.Errorf("Expected redis database 0, but got %v", config.RedisDatabase)
	}

	if config.RedisKeyPrefix != "probe:" {
		t.Errorf("Expected redis prefix probe:, but got %v", config.RedisKeyPrefix)
	}

	if config.RedisKeyTimeout != 2592000 {
		t.Errorf("Expected redis timeout 2592000, but got %v", config.RedisKeyTimeout)
	}
}
