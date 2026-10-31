package storage

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/petaki/probe/config"
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
	var errs []error

	if s.Config.DataLogEnabled {
		errs = append(errs, s.saveDataLog(m))
	}

	if s.Config.AlarmEnabled {
		errs = append(errs, s.saveAlarm(m))
	}

	if !s.Config.DataLogEnabled {
		errs = append(errs, s.printValue(m))
	}

	return errors.Join(errs...)
}

// IsPathIgnored function.
func (s *Storage) IsPathIgnored(path string) bool {
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
