package config

import (
	"errors"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestLoadAndParse(t *testing.T) {
	resetEnv(t)

	err := godotenv.Load("../.env.example")
	if err != nil {
		t.Errorf("Cannot load the .env.example file.")
	}

	config, err := Load()
	if err != nil {
		t.Errorf("Cannot load the environment variables.")
	}

	if config.Name != "probe" {
		t.Errorf("Expected probe name probe, but got %v", config.Name)
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

	if !config.DataLogEnabled {
		t.Errorf("Expected data log enabled true, but got %v", config.DataLogEnabled)
	}

	if config.DataLogTimeout != 2592000 {
		t.Errorf("Expected data log timeout 2592000, but got %v", config.DataLogTimeout)
	}

	if config.AlarmEnabled {
		t.Errorf("Expected alarm enabled false, but got %v", config.AlarmEnabled)
	}

	if config.AlarmCPUPercent != 0 {
		t.Errorf("Expected alarm cpu percent 0 (alarm disabled), but got %v", config.AlarmCPUPercent)
	}

	if config.AlarmMemoryPercent != 0 {
		t.Errorf("Expected alarm memory percent 0 (alarm disabled), but got %v", config.AlarmMemoryPercent)
	}

	if config.AlarmDiskPercent != 0 {
		t.Errorf("Expected alarm disk percent 0 (alarm disabled), but got %v", config.AlarmDiskPercent)
	}

	if config.AlarmLoadValue != 0 {
		t.Errorf("Expected alarm load value 0 (alarm disabled), but got %v", config.AlarmLoadValue)
	}

	if config.AlarmWebhookMethod != "" {
		t.Errorf("Expected alarm webhook method empty (alarm disabled), but got %v", config.AlarmWebhookMethod)
	}

	if config.AlarmWebhookURL != "" {
		t.Errorf("Expected alarm webhook URL empty (alarm disabled), but got %v", config.AlarmWebhookURL)
	}

	if config.AlarmWebhookHeader != nil {
		t.Errorf("Expected alarm webhook header nil (alarm disabled), but got %v", config.AlarmWebhookHeader)
	}

	if config.AlarmWebhookBody != "" {
		t.Errorf("Expected alarm webhook body empty (alarm disabled), but got %v", config.AlarmWebhookBody)
	}

	if config.AlarmFilterEnabled {
		t.Errorf("Expected alarm filter enabled false, but got %v", config.AlarmFilterEnabled)
	}

	if config.AlarmFilterWait != 0 {
		t.Errorf("Expected alarm filter wait 0 (filter disabled), but got %v", config.AlarmFilterWait)
	}

	if config.AlarmFilterSleep != 0 {
		t.Errorf("Expected alarm filter sleep 0 (filter disabled), but got %v", config.AlarmFilterSleep)
	}

	if config.LogTailEnabled {
		t.Errorf("Expected log tail enabled false, but got %v", config.LogTailEnabled)
	}

	if config.LogTailFiles != nil {
		t.Errorf("Expected log tail files nil (log tail disabled), but got %v", config.LogTailFiles)
	}

	if config.LogTailLines != 0 {
		t.Errorf("Expected log tail lines 0 (log tail disabled), but got %v", config.LogTailLines)
	}

	if config.LogTailBufferSize != 0 {
		t.Errorf("Expected log tail buffer size 0 (log tail disabled), but got %v", config.LogTailBufferSize)
	}

	if config.LogTailLimit != 0 {
		t.Errorf("Expected log tail limit 0 (log tail disabled), but got %v", config.LogTailLimit)
	}

	if config.LogTailTimeout != 0 {
		t.Errorf("Expected log tail timeout 0 (log tail disabled), but got %v", config.LogTailTimeout)
	}
}

func TestLoadAndParseWithAllFeaturesEnabled(t *testing.T) {
	resetEnv(t)

	err := godotenv.Load("../.env.example")
	if err != nil {
		t.Errorf("Cannot load the .env.example file.")
	}

	t.Setenv(envAlarmEnabled, "true")
	t.Setenv(envAlarmFilterEnabled, "true")
	t.Setenv(envLogTailEnabled, "true")

	config, err := Load()
	if err != nil {
		t.Fatalf("Cannot load the environment variables: %v", err)
	}

	if config.AlarmCPUPercent != 30 {
		t.Errorf("Expected alarm cpu percent 30, but got %v", config.AlarmCPUPercent)
	}

	if config.AlarmMemoryPercent != 50 {
		t.Errorf("Expected alarm memory percent 50, but got %v", config.AlarmMemoryPercent)
	}

	if config.AlarmDiskPercent != 80 {
		t.Errorf("Expected alarm disk percent 80, but got %v", config.AlarmDiskPercent)
	}

	if config.AlarmLoadValue != 1.0 {
		t.Errorf("Expected alarm load value 1.0, but got %v", config.AlarmLoadValue)
	}

	if config.AlarmWebhookMethod != "POST" {
		t.Errorf("Expected alarm webhook method POST, but got %v", config.AlarmWebhookMethod)
	}

	if config.AlarmWebhookURL != "http://127.0.0.1:4000/alarm" {
		t.Errorf("Expected alarm webhook URL http://127.0.0.1:4000/alarm, but got %v", config.AlarmWebhookURL)
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

	if config.AlarmWebhookBody != "{\"probe\": \"%p\", \"name\": \"%n\", \"alarm\": %a, \"used\": %u, \"timestamp_rfc3339\": \"%t\", \"timestamp_unix\": %x, \"link\": \"%l\"}" {
		t.Errorf("Expected alarm webhook body, but got %v", config.AlarmWebhookBody)
	}

	if config.AlarmFilterWait != 5 {
		t.Errorf("Expected alarm filter wait 5, but got %v", config.AlarmFilterWait)
	}

	if config.AlarmFilterSleep != 300 {
		t.Errorf("Expected alarm filter sleep 300, but got %v", config.AlarmFilterSleep)
	}

	for _, value := range config.LogTailFiles {
		if value == "/var/log/syslog" {
			continue
		}

		if value == "/var/log/auth.log" {
			continue
		}

		t.Errorf("Expected log tail files [/var/log/syslog /var/log/auth.log], but got %v", config.LogTailFiles)
	}

	if config.LogTailLines != 10 {
		t.Errorf("Expected log tail lines 10, but got %v", config.LogTailLines)
	}

	if config.LogTailBufferSize != 4096 {
		t.Errorf("Expected log tail buffer size 4096, but got %v", config.LogTailBufferSize)
	}

	if config.LogTailLimit != 60 {
		t.Errorf("Expected log tail limit 60, but got %v", config.LogTailLimit)
	}

	if config.LogTailTimeout != 172800 {
		t.Errorf("Expected log tail timeout 172800, but got %v", config.LogTailTimeout)
	}
}

func TestLoadAllowsMissingSubOptionsWhenFeatureDisabled(t *testing.T) {
	resetEnv(t)

	t.Setenv(envName, "probe")
	t.Setenv(envDiskIgnored, "/dev")
	t.Setenv(envDataLogEnabled, "false")
	t.Setenv(envAlarmEnabled, "false")
	t.Setenv(envAlarmFilterEnabled, "false")
	t.Setenv(envLogTailEnabled, "false")

	config, err := Load()
	if err != nil {
		t.Fatalf("Expected Load to succeed without sub-options, but got: %v", err)
	}

	if config.RedisURL != "" {
		t.Errorf("Expected redis URL empty (no redis features enabled), but got %v", config.RedisURL)
	}
}

func TestLoadRequiresProbeName(t *testing.T) {
	resetEnv(t)

	t.Setenv(envDiskIgnored, "/dev")
	t.Setenv(envDataLogEnabled, "false")
	t.Setenv(envAlarmEnabled, "false")
	t.Setenv(envAlarmFilterEnabled, "false")
	t.Setenv(envLogTailEnabled, "false")

	_, err := Load()
	if err == nil {
		t.Fatal("Expected Load to fail when probe name is missing, but got nil")
	}
}

func TestLoadRejectsEmptyProbeName(t *testing.T) {
	resetEnv(t)

	t.Setenv(envName, "")
	t.Setenv(envDiskIgnored, "/dev")
	t.Setenv(envDataLogEnabled, "false")
	t.Setenv(envAlarmEnabled, "false")
	t.Setenv(envAlarmFilterEnabled, "false")
	t.Setenv(envLogTailEnabled, "false")

	_, err := Load()
	if err == nil {
		t.Fatal("Expected Load to fail when probe name is empty, but got nil")
	}
}

func TestLoadRequiresRedisURLWhenDataLogEnabled(t *testing.T) {
	resetEnv(t)

	t.Setenv(envName, "probe")
	t.Setenv(envDiskIgnored, "/dev")
	t.Setenv(envDataLogEnabled, "true")
	t.Setenv(envDataLogTimeout, "60")
	t.Setenv(envAlarmEnabled, "false")
	t.Setenv(envAlarmFilterEnabled, "false")
	t.Setenv(envLogTailEnabled, "false")

	_, err := Load()
	if err == nil {
		t.Fatal("Expected Load to fail when data log is enabled but redis URL is missing, but got nil")
	}
}

func TestLoadRequiresRedisURLWhenAlarmFilterEnabled(t *testing.T) {
	resetEnv(t)

	t.Setenv(envName, "probe")
	t.Setenv(envDiskIgnored, "/dev")
	t.Setenv(envDataLogEnabled, "false")
	t.Setenv(envAlarmEnabled, "false")
	t.Setenv(envAlarmFilterEnabled, "true")
	t.Setenv(envAlarmFilterWait, "5")
	t.Setenv(envAlarmFilterSleep, "300")
	t.Setenv(envLogTailEnabled, "false")

	_, err := Load()
	if err == nil {
		t.Fatal("Expected Load to fail when alarm filter is enabled but redis URL is missing, but got nil")
	}
}

func TestLoadFailsWhenEnabledFeatureMissingSubOptions(t *testing.T) {
	cases := []struct {
		name string
		env  map[string]string
	}{
		{
			name: "data log missing timeout",
			env: map[string]string{
				envDataLogEnabled: "true",
				envRedisURL:       "redis://127.0.0.1:6379/0",
			},
		},
		{
			name: "alarm missing webhook url",
			env: map[string]string{
				envAlarmEnabled:       "true",
				envAlarmCPUPercent:    "30",
				envAlarmMemoryPercent: "50",
				envAlarmDiskPercent:   "80",
				envAlarmLoadValue:     "1.0",
				envAlarmWebhookMethod: "POST",
				envAlarmWebhookHeader: "{}",
				envAlarmWebhookBody:   "{}",
			},
		},
		{
			name: "alarm filter missing wait",
			env: map[string]string{
				envAlarmFilterEnabled: "true",
				envAlarmFilterSleep:   "300",
				envRedisURL:           "redis://127.0.0.1:6379/0",
			},
		},
		{
			name: "log tail missing files",
			env: map[string]string{
				envLogTailEnabled:    "true",
				envLogTailLines:      "10",
				envLogTailBufferSize: "4096",
				envLogTailLimit:      "60",
				envLogTailTimeout:    "172800",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resetEnv(t)

			t.Setenv(envName, "probe")
			t.Setenv(envDiskIgnored, "/dev")
			t.Setenv(envDataLogEnabled, "false")
			t.Setenv(envAlarmEnabled, "false")
			t.Setenv(envAlarmFilterEnabled, "false")
			t.Setenv(envLogTailEnabled, "false")

			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			_, err := Load()
			if err == nil {
				t.Fatalf("Expected Load to fail (%s), but got nil", tc.name)
			}
		})
	}
}

func resetEnv(t *testing.T) {
	t.Helper()

	keys := []string{
		envName,
		envDiskIgnored,
		envRedisURL,
		envDataLogEnabled,
		envDataLogTimeout,
		envAlarmEnabled,
		envAlarmCPUPercent,
		envAlarmMemoryPercent,
		envAlarmDiskPercent,
		envAlarmLoadValue,
		envAlarmWebhookMethod,
		envAlarmWebhookURL,
		envAlarmWebhookHeader,
		envAlarmWebhookBody,
		envAlarmFilterEnabled,
		envAlarmFilterWait,
		envAlarmFilterSleep,
		envLogTailEnabled,
		envLogTailFiles,
		envLogTailLines,
		envLogTailBufferSize,
		envLogTailLimit,
		envLogTailTimeout,
	}

	for _, key := range keys {
		original, hadKey := os.LookupEnv(key)
		os.Unsetenv(key)

		if hadKey {
			t.Cleanup(func() { os.Setenv(key, original) })
		} else {
			t.Cleanup(func() { os.Unsetenv(key) })
		}
	}
}

func TestLoadRejectsAlarmFilterWaitWithoutDataLog(t *testing.T) {
	resetEnv(t)

	t.Setenv(envName, "probe")
	t.Setenv(envDiskIgnored, "/dev")
	t.Setenv(envDataLogEnabled, "false")
	t.Setenv(envAlarmEnabled, "true")
	t.Setenv(envAlarmCPUPercent, "30")
	t.Setenv(envAlarmMemoryPercent, "50")
	t.Setenv(envAlarmDiskPercent, "80")
	t.Setenv(envAlarmLoadValue, "1.0")
	t.Setenv(envAlarmWebhookMethod, "POST")
	t.Setenv(envAlarmWebhookURL, "http://127.0.0.1:4000/alarm")
	t.Setenv(envAlarmWebhookHeader, "{}")
	t.Setenv(envAlarmWebhookBody, "{}")
	t.Setenv(envAlarmFilterEnabled, "true")
	t.Setenv(envAlarmFilterWait, "5")
	t.Setenv(envAlarmFilterSleep, "300")
	t.Setenv(envRedisURL, "redis://127.0.0.1:6379/0")
	t.Setenv(envLogTailEnabled, "false")

	_, err := Load()
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("Expected Load to fail with ErrInvalidValue, but got: %v", err)
	}
}

func TestLoadAllowsAlarmFilterWaitWithDataLog(t *testing.T) {
	resetEnv(t)

	t.Setenv(envName, "probe")
	t.Setenv(envDiskIgnored, "/dev")
	t.Setenv(envDataLogEnabled, "true")
	t.Setenv(envDataLogTimeout, "2592000")
	t.Setenv(envAlarmEnabled, "true")
	t.Setenv(envAlarmCPUPercent, "30")
	t.Setenv(envAlarmMemoryPercent, "50")
	t.Setenv(envAlarmDiskPercent, "80")
	t.Setenv(envAlarmLoadValue, "1.0")
	t.Setenv(envAlarmWebhookMethod, "POST")
	t.Setenv(envAlarmWebhookURL, "http://127.0.0.1:4000/alarm")
	t.Setenv(envAlarmWebhookHeader, "{}")
	t.Setenv(envAlarmWebhookBody, "{}")
	t.Setenv(envAlarmFilterEnabled, "true")
	t.Setenv(envAlarmFilterWait, "5")
	t.Setenv(envAlarmFilterSleep, "300")
	t.Setenv(envRedisURL, "redis://127.0.0.1:6379/0")
	t.Setenv(envLogTailEnabled, "false")

	_, err := Load()
	if err != nil {
		t.Fatalf("Expected Load to succeed with the data log enabled, but got: %v", err)
	}
}

func TestLoadAllowsAlarmFilterSleepOnlyWithoutDataLog(t *testing.T) {
	resetEnv(t)

	t.Setenv(envName, "probe")
	t.Setenv(envDiskIgnored, "/dev")
	t.Setenv(envDataLogEnabled, "false")
	t.Setenv(envAlarmEnabled, "true")
	t.Setenv(envAlarmCPUPercent, "30")
	t.Setenv(envAlarmMemoryPercent, "50")
	t.Setenv(envAlarmDiskPercent, "80")
	t.Setenv(envAlarmLoadValue, "1.0")
	t.Setenv(envAlarmWebhookMethod, "POST")
	t.Setenv(envAlarmWebhookURL, "http://127.0.0.1:4000/alarm")
	t.Setenv(envAlarmWebhookHeader, "{}")
	t.Setenv(envAlarmWebhookBody, "{}")
	t.Setenv(envAlarmFilterEnabled, "true")
	t.Setenv(envAlarmFilterWait, "0")
	t.Setenv(envAlarmFilterSleep, "300")
	t.Setenv(envRedisURL, "redis://127.0.0.1:6379/0")
	t.Setenv(envLogTailEnabled, "false")

	_, err := Load()
	if err != nil {
		t.Fatalf("Expected Load to succeed with the alarm filter wait disabled, but got: %v", err)
	}
}

func TestLoadRejectsEmptyLogTailFiles(t *testing.T) {
	resetEnv(t)

	t.Setenv(envName, "probe")
	t.Setenv(envDiskIgnored, "/dev")
	t.Setenv(envDataLogEnabled, "false")
	t.Setenv(envAlarmEnabled, "false")
	t.Setenv(envAlarmFilterEnabled, "false")
	t.Setenv(envLogTailEnabled, "true")
	t.Setenv(envLogTailFiles, "")
	t.Setenv(envLogTailLines, "10")
	t.Setenv(envLogTailBufferSize, "4096")
	t.Setenv(envLogTailLimit, "60")
	t.Setenv(envLogTailTimeout, "172800")

	_, err := Load()
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("Expected Load to fail with ErrInvalidValue, but got: %v", err)
	}
}

func TestLoadSkipsEmptyLogTailFileEntries(t *testing.T) {
	resetEnv(t)

	t.Setenv(envName, "probe")
	t.Setenv(envDiskIgnored, "/dev")
	t.Setenv(envDataLogEnabled, "false")
	t.Setenv(envAlarmEnabled, "false")
	t.Setenv(envAlarmFilterEnabled, "false")
	t.Setenv(envLogTailEnabled, "true")
	t.Setenv(envLogTailFiles, "/var/log/syslog,,/var/log/auth.log,")
	t.Setenv(envLogTailLines, "10")
	t.Setenv(envLogTailBufferSize, "4096")
	t.Setenv(envLogTailLimit, "60")
	t.Setenv(envLogTailTimeout, "172800")

	config, err := Load()
	if err != nil {
		t.Fatalf("Expected Load to succeed, but got: %v", err)
	}

	want := []string{"/var/log/syslog", "/var/log/auth.log"}
	if len(config.LogTailFiles) != len(want) {
		t.Fatalf("Expected %v, but got %v", want, config.LogTailFiles)
	}

	for i, value := range want {
		if config.LogTailFiles[i] != value {
			t.Errorf("Expected %v, but got %v", want, config.LogTailFiles)
		}
	}
}
