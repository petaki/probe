package storage

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/petaki/probe/config"
	"github.com/petaki/probe/model"
)

// Storage type.
type Storage struct {
	Prefix string
	Pool   *redis.Pool
}

// New function.
func New(config *config.Config) Storage {
	storage := Storage{}

	storage.Prefix = config.RedisPrefix
	storage.Pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			options := []redis.DialOption{
				redis.DialDatabase(config.RedisDatabase),
			}

			if config.RedisPassword != "" {
				options = append(options, redis.DialPassword(config.RedisPassword))
			}

			return redis.Dial("tcp", net.JoinHostPort(config.RedisHost, config.RedisPort), options...)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}

			_, err := c.Do("PING")

			return err
		},
	}

	return storage
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

	field := s.field()

	switch value := m.(type) {
	case model.CPU:
		_, err = s.hset(key, field, value.Used)
	case model.Disk:
		_, err = s.hset(key, field, value.Used)
	case model.Memory:
		_, err = s.hset(key, field, value.Used)
	}

	if err != nil {
		return err
	}

	if !exists {
		_, err := s.expire(key, s.timeout())
		if err != nil {
			return err
		}
	}

	return nil
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

	return date.String()
}

func (s *Storage) timeout() int {
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

	return int(7*24*time.Hour - now.Sub(date))
}

func (s *Storage) key(m interface{}) (string, error) {
	weekday := int(time.Now().Weekday())

	switch value := m.(type) {
	case model.CPU:
		return fmt.Sprintf("%scpu:%v", s.Prefix, weekday), nil
	case model.Memory:
		return fmt.Sprintf("%smemory:%v", s.Prefix, weekday), nil
	case model.Disk:
		encodedPath := base64.StdEncoding.EncodeToString([]byte(value.Path))

		return fmt.Sprintf("%sdisk:%v:%s", s.Prefix, weekday, encodedPath), nil
	}

	return "", errors.New("Unknown model type")
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
