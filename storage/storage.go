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
	Client *http.Client
	Config *config.Config
	Pool   *redis.Pool
}

// New function.
func New(config *config.Config) *Storage {
	var client *http.Client

	if config.AlarmEnabled {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	return &Storage{
		Client: client,
		Config: config,
		Pool: &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(config.RedisURL)
			},
		},
	}
}

// Save function.
func (s *Storage) Save(m interface{}) error {
	key, err := s.key(m)
	if err != nil {
		return err
	}

	exists, err := s.exists(key)
	if err != nil {
		return err
	}

	callAlarm := false

	switch value := m.(type) {
	case model.CPU:
		_, err = s.hset(key, s.field(), value.Used)
		callAlarm = s.Config.AlarmEnabled && s.Config.AlarmCPUPercent > 0 && value.Used >= s.Config.AlarmCPUPercent
	case model.Disk:
		_, err = s.hset(key, s.field(), value.Used)
		callAlarm = s.Config.AlarmEnabled && s.Config.AlarmMemoryPercent > 0 && value.Used >= s.Config.AlarmMemoryPercent
	case model.Memory:
		_, err = s.hset(key, s.field(), value.Used)
		callAlarm = s.Config.AlarmEnabled && s.Config.AlarmDiskPercent > 0 && value.Used >= s.Config.AlarmDiskPercent
	}

	if err != nil {
		return err
	}

	if !exists {
		_, err = s.expire(key, s.Config.RedisKeyTimeout)
		if err != nil {
			return err
		}
	}

	if callAlarm {
		key, err = s.alarmKey(m)
		if err != nil {
			return err
		}

		exists, err = s.exists(key)
		if err != nil {
			return err
		}

		if exists {
			return nil
		}

		err = s.callAlarm(m)
		if err != nil {
			return nil
		}

		_, err = s.set(key, "true")
		if err != nil {
			return err
		}

		_, err = s.expire(key, s.Config.AlarmTimeout)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) callAlarm(m interface{}) error {
	var name string
	var used float64

	switch value := m.(type) {
	case model.CPU:
		name = "CPU"
		used = value.Used
	case model.Memory:
		name = "Memory"
		used = value.Used
	case model.Disk:
		name = fmt.Sprintf("Disk:%s", value.Path)
		used = value.Used
	default:
		return ErrUnknownModelType
	}

	data := strings.ReplaceAll(s.Config.AlarmWebhookBody, "%n", name)
	data = strings.ReplaceAll(data, "%u", fmt.Sprintf("%.2f", used))

	req, err := http.NewRequest(s.Config.AlarmWebhookMethod, s.Config.AlarmWebhookURL, bytes.NewBuffer([]byte(data)))
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

func (s *Storage) alarmKey(m interface{}) (string, error) {
	switch value := m.(type) {
	case model.CPU:
		return fmt.Sprintf("%salarm:cpu", s.Config.RedisKeyPrefix), nil
	case model.Memory:
		return fmt.Sprintf("%salarm:memory", s.Config.RedisKeyPrefix), nil
	case model.Disk:
		encodedPath := base64.StdEncoding.EncodeToString([]byte(value.Path))

		return fmt.Sprintf("%salarm:disk:%s", s.Config.RedisKeyPrefix, encodedPath), nil
	default:
		return "", ErrUnknownModelType
	}
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

func (s *Storage) exists(key string) (bool, error) {
	conn := s.Pool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("EXISTS", key))
}

func (s *Storage) set(key string, value interface{}) (string, error) {
	conn := s.Pool.Get()
	defer conn.Close()

	return redis.String(conn.Do("SET", key, value))
}

func (s *Storage) hset(key string, field string, value interface{}) (bool, error) {
	conn := s.Pool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("HSET", key, field, value))
}

func (s *Storage) expire(key string, timeout int) (bool, error) {
	conn := s.Pool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("EXPIRE", key, timeout))
}
