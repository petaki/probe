package storage

import (
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

	cases := []struct {
		name string
		now  time.Time
	}{
		{"same day", time.Date(2026, 11, 20, 12, 1, 30, 0, time.Local)},
		{"across midnight", time.Date(2026, 11, 20, 0, 1, 30, 0, time.Local)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
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

			end := tc.now.Truncate(time.Minute).Add(-time.Minute)

			for minute := end.Add(-3 * time.Minute); !minute.After(end); minute = minute.Add(time.Minute) {
				day := time.Date(minute.Year(), minute.Month(), minute.Day(), 0, 0, 0, 0, minute.Location())
				key := fmt.Sprintf("%s:cpu:%d", cfg.Name, day.Unix())

				_, err := conn.Do("HSET", key, strconv.FormatInt(minute.Unix(), 10), 90)
				if err != nil {
					t.Fatal(err)
				}

				_, err = conn.Do("EXPIRE", key, 60)
				if err != nil {
					t.Fatal(err)
				}
			}

			err := storage.filterAlarm(model.CPU{Used: 90}, tc.now)
			if err != nil {
				t.Fatal(err)
			}

			if calls != 1 {
				t.Errorf("Expected 1 webhook call after %d minutes above the threshold, got %d", cfg.AlarmFilterWait, calls)
			}
		})
	}
}
