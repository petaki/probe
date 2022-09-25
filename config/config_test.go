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

	for _, value := range config.DiskIgnored {
		if value == "/dev" {
			continue
		}

		if value == "/var/lib/docker/*" {
			continue
		}

		t.Errorf("Expected disk ignored [/dev /var/lib/docker/*], but got %v", config.DiskIgnored)
	}

	if config.RedisURL != "redis://127.0.0.1:6379/0" {
		t.Errorf("Expected redis URL redis://127.0.0.1:6379/0, but got %v", config.RedisURL)
	}

	if config.RedisKeyPrefix != "probe:" {
		t.Errorf("Expected redis prefix probe:, but got %v", config.RedisKeyPrefix)
	}

	if !config.DataLogEnabled {
		t.Errorf("Expected data log enabled true, but got %v", config.DataLogEnabled)
	}

	if config.DataLogTimeout != 2592000 {
		t.Errorf("Expected data log timeout 2592000, but got %v", config.DataLogTimeout)
	}

	if config.AlarmEnabled {
		t.Errorf("Expected alarm enabled false, but got %v", config.AlarmEnabled)
	}

	if config.AlarmTimeout != 300 {
		t.Errorf("Expected alarm timeout 300, but got %v", config.AlarmTimeout)
	}

	if config.AlarmCPUPercent != 30 {
		t.Errorf("Expected alarm cpu percent 20, but got %v", config.AlarmCPUPercent)
	}

	if config.AlarmMemoryPercent != 50 {
		t.Errorf("Expected alarm memory percent 50, but got %v", config.AlarmMemoryPercent)
	}

	if config.AlarmDiskPercent != 80 {
		t.Errorf("Expected alarm disk percent 80, but got %v", config.AlarmDiskPercent)
	}

	if config.AlarmWebhookMethod != "POST" {
		t.Errorf("Expected alarm webhook method POST, but got %v", config.RedisKeyPrefix)
	}

	if config.AlarmWebhookURL != "http://127.0.0.1:4000/alarm" {
		t.Errorf("Expected alarm webhook URL http://127.0.0.1:4000/alarm, but got %v", config.RedisURL)
	}

	for name, value := range config.AlarmWebhookHeader {
		if name == "Accept" && value == "application/json" {
			continue
		}

		if name == "Authorization" && value == "Bearer TOKEN" {
			continue
		}

		t.Errorf("Expected alarm webhook header map[Accept:application/json Authorization:Bearer TOKEN], but got %v", config.AlarmWebhookHeader)
	}

	if config.AlarmWebhookBody != "{\"probe\": \"%p\", \"name\": \"%n\", \"alarm\": %a, \"used\": %u, \"timestamp\": \"%t\", \"link\": \"%l\"}" {
		t.Errorf("Expected alarm webhook body {\"probe\": \"%%p\", \"name\": \"%%n\", \"alarm\": %%a, \"used\": %%u, \"timestamp\": \"%%t\", \"link\": \"%%l\"}, but got %v", config.AlarmWebhookBody)
	}
}
