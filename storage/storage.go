package storage

import (
	"net/http"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/petaki/probe/config"
	"github.com/petaki/probe/model"
)

// Storage type.
type Storage struct {
	Config     *config.Config
	Pool       *redis.Pool
	Client     *http.Client
	normalizer *strings.Replacer
}

// New function.
func New(config *config.Config) *Storage {
	var pool *redis.Pool
	var client *http.Client

	if config.DataLogEnabled || config.AlarmFilterEnabled {
		pool = newPool(config)
	}

	if config.AlarmEnabled {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	return &Storage{
		Config:     config,
		Pool:       pool,
		Client:     client,
		normalizer: strings.NewReplacer(":", "_", "|", "_"),
	}
}

// Save function.
func (s *Storage) Save(m any) error {
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

func (s *Storage) isPathIgnored(path string) bool {
	for _, pattern := range s.Config.DiskIgnored {
		value, hasSuffix := strings.CutSuffix(pattern, "*")
		value, hasPrefix := strings.CutPrefix(value, "*")

		switch {
		case hasPrefix && hasSuffix:
			if strings.Contains(path, value) {
				return true
			}
		case hasPrefix:
			if strings.HasSuffix(path, value) {
				return true
			}
		case hasSuffix:
			if strings.HasPrefix(path, value) {
				return true
			}
		default:
			if value == path {
				return true
			}
		}
	}

	return false
}
