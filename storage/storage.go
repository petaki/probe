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

	if config.DataLogEnabled {
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

	switch value := m.(type) {
	case model.CPU:
		err = conn.Send(
			"HSET", key, s.field(), value.Used,
		)
	case model.Memory:
		err = conn.Send(
			"HSET", key, s.field(), value.Used,
		)
	case model.Disk:
		err = conn.Send(
			"HSET", key, s.field(), value.Used,
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

	if s.Config.DataLogEnabled {
		conn := s.Pool.Get()
		defer conn.Close()

		alarmKey, err := s.alarmKey(m)
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

		err = s.callAlarm(m)
		if err != nil {
			return err
		}

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
			"EXPIRE", alarmKey, s.Config.AlarmTimeout,
		)
		if err != nil {
			return err
		}

		_, err = conn.Do("EXEC")
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

func (s *Storage) callAlarm(m interface{}) error {
	var name string
	var used float64
	var alarm float64

	switch value := m.(type) {
	case model.CPU:
		name = "CPU"
		alarm = s.Config.AlarmCPUPercent
		used = value.Used
	case model.Memory:
		name = "Memory"
		alarm = s.Config.AlarmMemoryPercent
		used = value.Used
	case model.Disk:
		name = fmt.Sprintf("Disk:%s", value.Path)
		alarm = s.Config.AlarmDiskPercent
		used = value.Used
	default:
		return ErrUnknownModelType
	}

	body := strings.ReplaceAll(s.Config.AlarmWebhookBody, "%p", strings.ReplaceAll(s.Config.RedisKeyPrefix, ":", ""))
	body = strings.ReplaceAll(body, "%n", name)
	body = strings.ReplaceAll(body, "%a", fmt.Sprintf("%.2f", alarm))
	body = strings.ReplaceAll(body, "%u", fmt.Sprintf("%.2f", used))
	body = strings.ReplaceAll(body, "%t", time.Now().Format(time.RFC3339))

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

	if resp.StatusCode > 400 {
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

func (s *Storage) field() string {
	now := time.Now()
	date := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		0,
		0,
		now.Location(),
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
