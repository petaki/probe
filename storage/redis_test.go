package storage

import (
	"encoding/base64"
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
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

func TestSaveDataLogKeepsNewestLogEntries(t *testing.T) {
	redisURL := os.Getenv("PROBE_TEST_REDIS_URL")
	if redisURL == "" {
		t.Skip("PROBE_TEST_REDIS_URL is not set")
	}

	cases := []struct {
		name   string
		seeded []string
		limit  int
		want   int
	}{
		{"under the limit", []string{"1000000001"}, 3, 2},
		{"at the limit", []string{"1000000001", "1000000002"}, 3, 3},
		{"over the limit", []string{"1000000001", "1000000002", "1000000003", "1000000004"}, 2, 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.Config{
				Name:           fmt.Sprintf("probe-test-%d", time.Now().UnixNano()),
				RedisURL:       redisURL,
				DataLogEnabled: true,
				LogTailLimit:   tc.limit,
				LogTailTimeout: 60,
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
			key := fmt.Sprintf("%s:log:%d:%s", cfg.Name, day.Unix(), base64.StdEncoding.EncodeToString([]byte("/var/log/app.log")))

			for _, field := range tc.seeded {
				_, err := conn.Do("HSET", key, field, "old")
				if err != nil {
					t.Fatal(err)
				}
			}

			_, err := conn.Do("EXPIRE", key, 60)
			if err != nil {
				t.Fatal(err)
			}

			err = storage.saveDataLog(model.Log{Path: "/var/log/app.log", Content: "new"})
			if err != nil {
				t.Fatal(err)
			}

			if time.Now().Day() != now.Day() {
				t.Skip("The day changed during the test")
			}

			fields, err := redis.StringMap(conn.Do("HGETALL", key))
			if err != nil {
				t.Fatal(err)
			}

			if len(fields) != tc.want {
				t.Fatalf("Expected %d entries, got %d: %v", tc.want, len(fields), fields)
			}

			newest := slices.Sorted(maps.Keys(fields))[len(fields)-1]
			if fields[newest] != "new" {
				t.Errorf("Expected the newest entry to hold the latest tail, got %q", fields[newest])
			}

			seededKept := tc.want - 1
			for _, field := range tc.seeded[:len(tc.seeded)-seededKept] {
				_, kept := fields[field]
				if kept {
					t.Errorf("Expected old entry %s to be trimmed", field)
				}
			}
		})
	}
}
