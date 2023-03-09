package storage

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/petaki/probe/config"
	"github.com/petaki/probe/model"
)

// Storage type.
type Storage struct {
	Config *config.Config
	Pool   *redis.Pool
	Client *http.Client
}

// New function.
func New(config *config.Config) *Storage {
	var pool *redis.Pool
	var client *http.Client

	if config.DataLogEnabled || config.AlarmFilterEnabled {
		pool = &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(config.RedisURL)
			},
		}
	}

	if config.AlarmEnabled {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	return &Storage{
		Config: config,
		Pool:   pool,
		Client: client,
	}
}

// Save function.
func (s *Storage) Save(m interface{}) error {
	var err error

	switch value := m.(type) {
	case model.Disk:
		if s.isPathIgnored(value.Path) {
			return nil
		}
	}

	if s.Config.DataLogEnabled {
		err = s.saveDataLog(m)
		if err != nil {
			return err
		}
	}

	if s.Config.AlarmEnabled {
		err = s.saveAlarm(m)
		if err != nil {
			return err
		}
	}

	if !s.Config.DataLogEnabled {
		err = s.printValue(m)
		if err != nil {
			return err
		}
	}

	return nil
}

// SaveAlarmConfig function.
func (s *Storage) SaveAlarmConfig() error {
	if !s.Config.DataLogEnabled {
		return nil
	}

	conn := s.Pool.Get()
	defer conn.Close()

	alarm := &model.Alarm{
		CPU:    s.Config.AlarmCPUPercent,
		Memory: s.Config.AlarmMemoryPercent,
		Disk:   s.Config.AlarmDiskPercent,
	}

	_, err := conn.Do(
		"HSET", redis.Args{}.Add(fmt.Sprintf("%salarm", s.Config.RedisKeyPrefix)).AddFlat(alarm)...,
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAlarmConfig function.
func (s *Storage) DeleteAlarmConfig() error {
	if !s.Config.DataLogEnabled {
		return nil
	}

	conn := s.Pool.Get()
	defer conn.Close()

	_, err := conn.Do(
		"DEL", fmt.Sprintf("%salarm", s.Config.RedisKeyPrefix),
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) saveDataLog(m interface{}) error {
	conn := s.Pool.Get()
	defer conn.Close()

	key, err := s.key(m)
	if err != nil {
		return err
	}

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return err
	}

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	now := time.Now()

	switch value := m.(type) {
	case model.CPU:
		err = conn.Send(
			"HSET", key, s.field(&now), value.Used,
		)
	case model.Memory:
		err = conn.Send(
			"HSET", key, s.field(&now), value.Used,
		)
	case model.Disk:
		err = conn.Send(
			"HSET", key, s.field(&now), value.Used,
		)
	}
	if err != nil {
		return err
	}

	if !exists {
		err = conn.Send(
			"EXPIRE", key, s.Config.DataLogTimeout,
		)
		if err != nil {
			return err
		}
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) saveAlarm(m interface{}) error {
	callAlarm := false

	switch value := m.(type) {
	case model.CPU:
		callAlarm = s.Config.AlarmCPUPercent > 0 && value.Used >= s.Config.AlarmCPUPercent
	case model.Memory:
		callAlarm = s.Config.AlarmMemoryPercent > 0 && value.Used >= s.Config.AlarmMemoryPercent
	case model.Disk:
		callAlarm = s.Config.AlarmDiskPercent > 0 && value.Used >= s.Config.AlarmDiskPercent
	default:
		return ErrUnknownModelType
	}

	if !callAlarm {
		return nil
	}

	if s.Config.AlarmFilterEnabled {
		err := s.filterAlarm(m)
		if err != nil {
			return err
		}

		return nil
	}

	err := s.callAlarm(m)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) filterAlarm(m interface{}) error {
	conn := s.Pool.Get()
	defer conn.Close()

	var alarmKey string
	var err error

	if s.Config.AlarmFilterSleep > 0 {
		alarmKey, err = s.alarmKey(m)
		if err != nil {
			return err
		}

		exists, err := redis.Bool(conn.Do("EXISTS", alarmKey))
		if err != nil {
			return err
		}

		if exists {
			return nil
		}
	}

	if s.Config.AlarmFilterWait > 1 {
		key, err := s.key(m)
		if err != nil {
			return err
		}

		var fields []string

		now := time.Now()
		end := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			now.Minute()-1,
			0,
			0,
			now.Location(),
		)

		start := end.Add(-time.Duration(s.Config.AlarmFilterWait-2) * time.Minute)

		for current := start; !current.After(end); current = current.Add(time.Minute) {
			fields = append(fields, s.field(&current))
		}

		values, err := redis.Float64s(conn.Do("HMGET", redis.Args{}.Add(key).AddFlat(fields)...))
		if err != nil {
			return err
		}

		for _, value := range values {
			switch m.(type) {
			case model.CPU:
				if value < s.Config.AlarmCPUPercent {
					return nil
				}
			case model.Memory:
				if value < s.Config.AlarmMemoryPercent {
					return nil
				}
			case model.Disk:
				if value < s.Config.AlarmDiskPercent {
					return nil
				}
			default:
				return ErrUnknownModelType
			}
		}
	}

	err = s.callAlarm(m)
	if err != nil {
		return err
	}

	if s.Config.AlarmFilterSleep > 0 {
		err = conn.Send("MULTI")
		if err != nil {
			return err
		}

		err = conn.Send(
			"SET", alarmKey, true,
		)
		if err != nil {
			return err
		}

		err = conn.Send(
			"EXPIRE", alarmKey, s.Config.AlarmFilterSleep,
		)
		if err != nil {
			return err
		}

		_, err = conn.Do("EXEC")
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) callAlarm(m interface{}) error {
	probe := strings.ReplaceAll(s.Config.RedisKeyPrefix, ":", "")

	var name string
	var used float64
	var alarm float64
	var link string

	switch value := m.(type) {
	case model.CPU:
		name = "CPU"
		alarm = s.Config.AlarmCPUPercent
		used = value.Used
		link = fmt.Sprintf("/cpu?probe=%s", probe)
	case model.Memory:
		name = "Memory"
		alarm = s.Config.AlarmMemoryPercent
		used = value.Used
		link = fmt.Sprintf("/memory?probe=%s", probe)
	case model.Disk:
		name = fmt.Sprintf("Disk:%s", value.Path)
		alarm = s.Config.AlarmDiskPercent
		used = value.Used
		link = fmt.Sprintf("/disk?probe=%s&path=%s", probe, value.Path)
	default:
		return ErrUnknownModelType
	}

	now := time.Now()

	body := strings.ReplaceAll(s.Config.AlarmWebhookBody, "%p", probe)
	body = strings.ReplaceAll(body, "%n", name)
	body = strings.ReplaceAll(body, "%a", fmt.Sprintf("%.2f", alarm))
	body = strings.ReplaceAll(body, "%u", fmt.Sprintf("%.2f", used))
	body = strings.ReplaceAll(body, "%t", now.Format(time.RFC3339))
	body = strings.ReplaceAll(body, "%x", strconv.FormatInt(now.Unix(), 10))
	body = strings.ReplaceAll(body, "%l", link)

	req, err := http.NewRequest(s.Config.AlarmWebhookMethod, s.Config.AlarmWebhookURL, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}

	for key, value := range s.Config.AlarmWebhookHeader {
		req.Header.Set(key, value)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return ErrBadStatusCode
	}

	return nil
}

func (s *Storage) printValue(m interface{}) error {
	switch value := m.(type) {
	case model.CPU:
		fmt.Printf("  ⚡ CPU: %.2f%%\n", value.Used)
	case model.Memory:
		fmt.Printf("  📦 Memory: %.2f%%\n", value.Used)
	case model.Disk:
		fmt.Printf("  💾 Disk:%s: %.2f%%\n", value.Path, value.Used)
	default:
		return ErrUnknownModelType
	}

	fmt.Println()

	return nil
}

func (s *Storage) isPathIgnored(path string) bool {
	for _, pattern := range s.Config.DiskIgnored {
		value := strings.ReplaceAll(pattern, "*", "")

		if pattern[0:1] == "*" && pattern[len(pattern)-1:] == "*" {
			if strings.Contains(path, value) {
				return true
			}

			continue
		}

		if pattern[0:1] == "*" {
			if strings.HasSuffix(path, value) {
				return true
			}

			continue
		}

		if pattern[len(pattern)-1:] == "*" {
			if strings.HasPrefix(path, value) {
				return true
			}

			continue
		}

		if value == path {
			return true
		}
	}

	return false
}

func (s *Storage) key(m interface{}) (string, error) {
	switch value := m.(type) {
	case model.CPU:
		return fmt.Sprintf("%scpu:%s", s.Config.RedisKeyPrefix, s.timestamp()), nil
	case model.Memory:
		return fmt.Sprintf("%smemory:%s", s.Config.RedisKeyPrefix, s.timestamp()), nil
	case model.Disk:
		encodedPath := base64.StdEncoding.EncodeToString([]byte(value.Path))

		return fmt.Sprintf("%sdisk:%s:%s", s.Config.RedisKeyPrefix, s.timestamp(), encodedPath), nil
	default:
		return "", ErrUnknownModelType
	}
}

func (s *Storage) timestamp() string {
	now := time.Now()
	date := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0,
		0,
		0,
		0,
		now.Location(),
	)

	return strconv.FormatInt(date.Unix(), 10)
}

func (s *Storage) field(t *time.Time) string {
	date := time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		0,
		0,
		t.Location(),
	)

	return strconv.FormatInt(date.Unix(), 10)
}

func (s *Storage) alarmKey(m interface{}) (string, error) {
	switch value := m.(type) {
	case model.CPU:
		return fmt.Sprintf("%salarm:cpu", s.Config.RedisKeyPrefix), nil
	case model.Memory:
		return fmt.Sprintf("%salarm:memory", s.Config.RedisKeyPrefix), nil
	case model.Disk:
		encodedPath := base64.StdEncoding.EncodeToString([]byte(value.Path))

		return fmt.Sprintf("%salarm:disk:%s", s.Config.RedisKeyPrefix, encodedPath), nil
	}

	return "", ErrUnknownModelType
}
