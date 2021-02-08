package storage

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/petaki/probe/config"
	"github.com/petaki/probe/model"
)

// Storage type.
type Storage struct {
	Config *config.Config
	Pool   *redis.Pool
}

// New function.
func New(config *config.Config) Storage {
	return Storage{
		Config: config,
		Pool: &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(config.RedisUrl)
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

	switch value := m.(type) {
	case model.CPU:
		_, err = s.hset(key, s.field(), value.Used)
	case model.Disk:
		_, err = s.hset(key, s.field(), value.Used)
	case model.Memory:
		_, err = s.hset(key, s.field(), value.Used)
	}

	if err != nil {
		return err
	}

	if !exists {
		_, err := s.expire(key, s.Config.RedisKeyTimeout)
		if err != nil {
			return err
		}
	}

	return nil
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
		return "", errors.New("Unknown model type")
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
