package storage

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/petaki/probe/config"
	"github.com/petaki/probe/model"
)

func TestFilterAlarmChecksEveryMinuteOfTheWait(t *testing.T) {
	redisURL := os.Getenv("PROBE_TEST_REDIS_URL")
	if redisURL == "" {
		t.Skip("PROBE_TEST_REDIS_URL is not set")
	}

	times := []struct {
		name string
		now  time.Time
	}{
		{"same day", time.Date(2026, 11, 20, 12, 1, 30, 0, time.Local)},
		{"across midnight", time.Date(2026, 11, 20, 0, 1, 30, 0, time.Local)},
	}

	watchers := []struct {
		name  string
		model any
		key   func(name string, day int64) string
		above string
		below string
	}{
		{
			"cpu",
			model.CPU{Used: 90},
			func(name string, day int64) string { return fmt.Sprintf("%s:cpu:%d", name, day) },
			"90",
			"10",
		},
		{
			"memory",
			model.Memory{Used: 90},
			func(name string, day int64) string { return fmt.Sprintf("%s:memory:%d", name, day) },
			"90",
			"10",
		},
		{
			"load",
			model.Load{Load1: 5, Load5: 5, Load15: 5},
			func(name string, day int64) string { return fmt.Sprintf("%s:load:%d", name, day) },
			"5.000000:5.000000:5.000000",
			"0.100000:0.100000:0.100000",
		},
		{
			"disk",
			model.Disk{Path: "/data", Used: 90},
			func(name string, day int64) string {
				return fmt.Sprintf("%s:disk:%d:%s", name, day, base64.StdEncoding.EncodeToString([]byte("/data")))
			},
			"90",
			"10",
		},
	}

	for _, tt := range times {
		for _, tw := range watchers {
			for _, dip := range []bool{false, true} {
				name := fmt.Sprintf("%s/%s/above every minute", tt.name, tw.name)
				want := 1

				if dip {
					name = fmt.Sprintf("%s/%s/one minute below", tt.name, tw.name)
					want = 0
				}

				t.Run(name, func(t *testing.T) {
					calls := 0

					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						calls++
					}))
					defer server.Close()

					cfg := &config.Config{
						Name:               fmt.Sprintf("probe-test-%d", time.Now().UnixNano()),
						RedisURL:           redisURL,
						AlarmEnabled:       true,
						AlarmCPUPercent:    30,
						AlarmMemoryPercent: 30,
						AlarmDiskPercent:   30,
						AlarmLoadValue:     1,
						AlarmWebhookMethod: http.MethodPost,
						AlarmWebhookURL:    server.URL,
						AlarmFilterEnabled: true,
						AlarmFilterWait:    5,
					}

					storage := &Storage{
						Config: cfg,
						Pool:   newPool(cfg),
						Client: server.Client(),
					}
					defer storage.Pool.Close()

					conn := storage.Pool.Get()
					defer conn.Close()

					end := tt.now.Truncate(time.Minute).Add(-time.Minute)
					start := end.Add(-3 * time.Minute)

					for minute := start; !minute.After(end); minute = minute.Add(time.Minute) {
						day := time.Date(minute.Year(), minute.Month(), minute.Day(), 0, 0, 0, 0, minute.Location())
						key := tw.key(cfg.Name, day.Unix())

						value := tw.above
						if dip && minute.Equal(start.Add(time.Minute)) {
							value = tw.below
						}

						_, err := conn.Do("HSET", key, strconv.FormatInt(minute.Unix(), 10), value)
						if err != nil {
							t.Fatal(err)
						}

						_, err = conn.Do("EXPIRE", key, 60)
						if err != nil {
							t.Fatal(err)
						}
					}

					err := storage.filterAlarm(tw.model, tt.now)
					if err != nil {
						t.Fatal(err)
					}

					if calls != want {
						t.Errorf("Expected %d webhook calls, got %d", want, calls)
					}
				})
			}
		}
	}
}
