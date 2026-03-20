package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	envDiskIgnored        = "PROBE_DISK_IGNORED"
	envRedisURL           = "PROBE_REDIS_URL"
	envRedisKeyPrefix     = "PROBE_REDIS_KEY_PREFIX"
	envDataLogEnabled     = "PROBE_DATA_LOG_ENABLED"
	envDataLogTimeout     = "PROBE_DATA_LOG_TIMEOUT"
	envAlarmEnabled       = "PROBE_ALARM_ENABLED"
	envAlarmCPUPercent    = "PROBE_ALARM_CPU_PERCENT"
	envAlarmMemoryPercent = "PROBE_ALARM_MEMORY_PERCENT"
	envAlarmDiskPercent   = "PROBE_ALARM_DISK_PERCENT"
	envAlarmLoadValue     = "PROBE_ALARM_LOAD_VALUE"
	envAlarmWebhookMethod = "PROBE_ALARM_WEBHOOK_METHOD"
	envAlarmWebhookURL    = "PROBE_ALARM_WEBHOOK_URL"
	envAlarmWebhookHeader = "PROBE_ALARM_WEBHOOK_HEADER"
	envAlarmWebhookBody   = "PROBE_ALARM_WEBHOOK_BODY"
	envAlarmFilterEnabled = "PROBE_ALARM_FILTER_ENABLED"
	envAlarmFilterWait    = "PROBE_ALARM_FILTER_WAIT"
	envAlarmFilterSleep   = "PROBE_ALARM_FILTER_SLEEP"
	envLogTailEnabled     = "PROBE_LOG_TAIL_ENABLED"
	envLogTailFiles       = "PROBE_LOG_TAIL_FILES"
	envLogTailLines       = "PROBE_LOG_TAIL_LINES"
	envLogTailBufferSize  = "PROBE_LOG_TAIL_BUFFER_SIZE"
	envLogTailLimit       = "PROBE_LOG_TAIL_LIMIT"
	envLogTailTimeout     = "PROBE_LOG_TAIL_TIMEOUT"
)

var envKeys = []string{
	envDiskIgnored,
	envRedisURL,
	envRedisKeyPrefix,
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

// Config type.
type Config struct {
	DiskIgnored        []string
	RedisURL           string
	RedisKeyPrefix     string
	DataLogEnabled     bool
	DataLogTimeout     int
	AlarmEnabled       bool
	AlarmCPUPercent    float64
	AlarmMemoryPercent float64
	AlarmDiskPercent   float64
	AlarmLoadValue     float64
	AlarmWebhookMethod string
	AlarmWebhookURL    string
	AlarmWebhookHeader map[string]string
	AlarmWebhookBody   string
	AlarmFilterEnabled bool
	AlarmFilterWait    int
	AlarmFilterSleep   int
	LogTailEnabled     bool
	LogTailFiles       []string
	LogTailLines       int
	LogTailBufferSize  int
	LogTailLimit       int
	LogTailTimeout     int
}

// Load function.
func Load() (*Config, error) {
	config := Config{}

	for _, key := range envKeys {
		value, hasKey := os.LookupEnv(key)
		if !hasKey {
			return nil, fmt.Errorf("%v is not defined", key)
		}

		err := config.parse(key, value)
		if err != nil {
			return nil, err
		}
	}

	return &config, nil
}

func (c *Config) parse(key string, value string) error {
	switch key {
	case envDiskIgnored:
		c.DiskIgnored = strings.Split(value, ",")
	case envRedisURL:
		c.RedisURL = value
	case envRedisKeyPrefix:
		c.RedisKeyPrefix = value
	case envDataLogEnabled:
		dataLogEnabled, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}

		c.DataLogEnabled = dataLogEnabled
	case envDataLogTimeout:
		dataLogTimeout, err := strconv.Atoi(value)
		if err != nil {
			return err
		}

		if dataLogTimeout < 1 {
			return ErrInvalidTimeout
		}

		c.DataLogTimeout = dataLogTimeout
	case envAlarmEnabled:
		alarmEnabled, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}

		c.AlarmEnabled = alarmEnabled
	case envAlarmCPUPercent:
		alarmCPUPercent, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}

		if alarmCPUPercent < 0 || alarmCPUPercent > 100 {
			return ErrInvalidPercent
		}

		c.AlarmCPUPercent = alarmCPUPercent
	case envAlarmMemoryPercent:
		alarmMemoryPercent, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}

		if alarmMemoryPercent < 0 || alarmMemoryPercent > 100 {
			return ErrInvalidPercent
		}

		c.AlarmMemoryPercent = alarmMemoryPercent
	case envAlarmDiskPercent:
		alarmDiskPercent, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}

		if alarmDiskPercent < 0 || alarmDiskPercent > 100 {
			return ErrInvalidPercent
		}

		c.AlarmDiskPercent = alarmDiskPercent
	case envAlarmLoadValue:
		alarmLoadValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}

		if alarmLoadValue < 0 {
			return ErrInvalidValue
		}

		c.AlarmLoadValue = alarmLoadValue
	case envAlarmWebhookMethod:
		c.AlarmWebhookMethod = value
	case envAlarmWebhookURL:
		c.AlarmWebhookURL = value
	case envAlarmWebhookHeader:
		var alarmWebhookHeader map[string]string

		err := json.Unmarshal([]byte(value), &alarmWebhookHeader)
		if err != nil {
			return err
		}

		c.AlarmWebhookHeader = alarmWebhookHeader
	case envAlarmWebhookBody:
		c.AlarmWebhookBody = value
	case envAlarmFilterEnabled:
		alarmFilterEnabled, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}

		c.AlarmFilterEnabled = alarmFilterEnabled
	case envAlarmFilterWait:
		alarmFilterWait, err := strconv.Atoi(value)
		if err != nil {
			return err
		}

		if alarmFilterWait < 0 {
			return ErrInvalidTimeout
		}

		c.AlarmFilterWait = alarmFilterWait
	case envAlarmFilterSleep:
		alarmFilterSleep, err := strconv.Atoi(value)
		if err != nil {
			return err
		}

		if alarmFilterSleep < 0 {
			return ErrInvalidTimeout
		}

		c.AlarmFilterSleep = alarmFilterSleep
	case envLogTailEnabled:
		logTailEnabled, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}

		c.LogTailEnabled = logTailEnabled
	case envLogTailFiles:
		c.LogTailFiles = strings.Split(value, ",")
	case envLogTailLines:
		logTailLines, err := strconv.Atoi(value)
		if err != nil {
			return err
		}

		if logTailLines < 1 {
			return ErrInvalidValue
		}

		c.LogTailLines = logTailLines
	case envLogTailBufferSize:
		logTailBufferSize, err := strconv.Atoi(value)
		if err != nil {
			return err
		}

		if logTailBufferSize < 1 {
			return ErrInvalidValue
		}

		c.LogTailBufferSize = logTailBufferSize
	case envLogTailLimit:
		logTailLimit, err := strconv.Atoi(value)
		if err != nil {
			return err
		}

		if logTailLimit < 1 {
			return ErrInvalidValue
		}

		c.LogTailLimit = logTailLimit
	case envLogTailTimeout:
		logTailTimeout, err := strconv.Atoi(value)
		if err != nil {
			return err
		}

		if logTailTimeout < 1 {
			return ErrInvalidTimeout
		}

		c.LogTailTimeout = logTailTimeout
	}

	return nil
}
