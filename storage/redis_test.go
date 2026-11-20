package storage

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/petaki/probe/config"
	"github.com/petaki/probe/model"
)

func TestSaveDataLogReportsFailedTransactionCommands(t *testing.T) {
	redisURL := os.Getenv("PROBE_TEST_REDIS_URL")
	if redisURL == "" {
		t.Skip("PROBE_TEST_REDIS_URL is not set")
	}

	cfg := &config.Config{
		Name:           fmt.Sprintf("probe-test-%d", time.Now().UnixNano()),
		RedisURL:       redisURL,
		DataLogEnabled: true,
		DataLogTimeout: 60,
	}

	storage := &Storage{
		Config: cfg,
		Pool:   newPool(cfg),
	}
	defer storage.Pool.Close()

	conn := storage.Pool.Get()
	defer conn.Close()

	now := time.Now()
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	key := fmt.Sprintf("%s:cpu:%d", cfg.Name, day.Unix())

	_, err := conn.Do("SET", key, "not a hash", "EX", 60)
	if err != nil {
		t.Fatal(err)
	}

	err = storage.saveDataLog(model.CPU{Used: 42})

	if time.Now().Day() != now.Day() {
		t.Skip("The day changed during the test")
	}

	if !errors.Is(err, ErrTransactionFailed) {
		t.Fatalf("Expected ErrTransactionFailed, got %v", err)
	}
}
