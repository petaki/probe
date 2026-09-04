package storage

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/petaki/probe/config"
	"github.com/petaki/probe/model"
)

func newPool(config *config.Config) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(
				config.RedisURL,
				redis.DialConnectTimeout(5*time.Second),
				redis.DialReadTimeout(5*time.Second),
				redis.DialWriteTimeout(5*time.Second),
			)
		},
	}
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
		Load:   s.Config.AlarmLoadValue,
	}

	_, err := conn.Do(
		"HSET", redis.Args{}.Add(fmt.Sprintf("%s:alarm", s.Config.Name)).AddFlat(alarm)...,
	)

	return err
}

// DeleteAlarmConfig function.
func (s *Storage) DeleteAlarmConfig() error {
	if !s.Config.DataLogEnabled {
		return nil
	}

	conn := s.Pool.Get()
	defer conn.Close()

	_, err := conn.Do(
		"DEL", fmt.Sprintf("%s:alarm", s.Config.Name),
	)

	return err
}

func (s *Storage) saveDataLog(m any) error {
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
	case []model.ProcessCPU:
		var v []string

		for _, p := range value {
			v = append(v, fmt.Sprintf("%s:%f", s.normalizer.Replace(p.Name), p.Used))
		}

		err = conn.Send(
			"HSET", key, s.field(&now), strings.Join(v, "|"),
		)
	case []model.ProcessMemory:
		var v []string

		for _, p := range value {
			v = append(v, fmt.Sprintf("%s:%f", s.normalizer.Replace(p.Name), p.Used))
		}

		err = conn.Send(
			"HSET", key, s.field(&now), strings.Join(v, "|"),
		)
	case model.Load:
		err = conn.Send(
			"HSET", key, s.field(&now), fmt.Sprintf("%f:%f:%f", value.Load1, value.Load5, value.Load15),
		)
	case model.Disk:
		err = conn.Send(
			"HSET", key, s.field(&now), value.Used,
		)
	case model.Log:
		err = conn.Send(
			"HSET", key, s.field(&now), value.Content,
		)
	}
	if err != nil {
		return err
	}

	if !exists {
		var timeout int
		if _, ok := m.(model.Log); ok {
			timeout = s.Config.LogTailTimeout
		} else {
			timeout = s.Config.DataLogTimeout
		}

		err = conn.Send("EXPIRE", key, timeout)
		if err != nil {
			return err
		}
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	_, ok := m.(model.Log)
	if ok && s.Config.LogTailLimit > 0 {
		return s.trimLogEntries(conn, key)
	}

	return nil
}

func (s *Storage) trimLogEntries(conn redis.Conn, key string) error {
	length, err := redis.Int(conn.Do("HLEN", key))
	if err != nil {
		return err
	}

	if length <= s.Config.LogTailLimit {
		return nil
	}

	fields, err := redis.Strings(conn.Do("HKEYS", key))
	if err != nil {
		return err
	}

	if len(fields) <= s.Config.LogTailLimit {
		return nil
	}

	sort.Strings(fields)

	remove := fields[:len(fields)-s.Config.LogTailLimit]

	_, err = conn.Do("HDEL", redis.Args{}.Add(key).AddFlat(remove)...)

	return err
}

func (s *Storage) saveAlarm(m any) error {
	callAlarm := false

	switch value := m.(type) {
	case model.CPU:
		callAlarm = s.Config.AlarmCPUPercent > 0 && value.Used >= s.Config.AlarmCPUPercent
	case model.Memory:
		callAlarm = s.Config.AlarmMemoryPercent > 0 && value.Used >= s.Config.AlarmMemoryPercent
	case []model.ProcessCPU:
		return nil
	case []model.ProcessMemory:
		return nil
	case model.Load:
		callAlarm = s.Config.AlarmLoadValue > 0 && (value.Load1 >= s.Config.AlarmLoadValue || value.Load5 >= s.Config.AlarmLoadValue || value.Load15 >= s.Config.AlarmLoadValue)
	case model.Disk:
		callAlarm = s.Config.AlarmDiskPercent > 0 && value.Used >= s.Config.AlarmDiskPercent
	case model.Log:
		return nil
	default:
		return ErrUnknownModelType
	}

	if !callAlarm {
		return nil
	}

	if s.Config.AlarmFilterEnabled {
		return s.filterAlarm(m)
	}

	return s.callAlarm(m)
}

func (s *Storage) filterAlarm(m any) error {
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

		switch m.(type) {
		case model.CPU, model.Memory, model.Disk:
			var threshold float64

			switch m.(type) {
			case model.CPU:
				threshold = s.Config.AlarmCPUPercent
			case model.Memory:
				threshold = s.Config.AlarmMemoryPercent
			case model.Disk:
				threshold = s.Config.AlarmDiskPercent
			}

			values, err := redis.Float64s(conn.Do("HMGET", redis.Args{}.Add(key).AddFlat(fields)...))
			if err != nil {
				return err
			}

			for _, value := range values {
				if value < threshold {
					return nil
				}
			}
		case model.Load:
			values, err := redis.Strings(conn.Do("HMGET", redis.Args{}.Add(key).AddFlat(fields)...))
			if err != nil {
				return err
			}

			for _, raw := range values {
				value := true
				segments := strings.SplitN(raw, ":", 3)

				if len(segments) != 3 {
					continue
				}

				for _, segment := range segments {
					segmentValue, err := strconv.ParseFloat(segment, 64)
					if err != nil {
						return err
					}

					value = value && segmentValue < s.Config.AlarmLoadValue
				}

				if value {
					return nil
				}
			}
		default:
			return ErrUnknownModelType
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

func (s *Storage) key(m any) (string, error) {
	switch value := m.(type) {
	case model.CPU:
		return fmt.Sprintf("%s:cpu:%s", s.Config.Name, s.timestamp()), nil
	case model.Memory:
		return fmt.Sprintf("%s:memory:%s", s.Config.Name, s.timestamp()), nil
	case []model.ProcessCPU:
		return fmt.Sprintf("%s:process:cpu:%s", s.Config.Name, s.timestamp()), nil
	case []model.ProcessMemory:
		return fmt.Sprintf("%s:process:memory:%s", s.Config.Name, s.timestamp()), nil
	case model.Load:
		return fmt.Sprintf("%s:load:%s", s.Config.Name, s.timestamp()), nil
	case model.Disk:
		encodedPath := base64.StdEncoding.EncodeToString([]byte(value.Path))

		return fmt.Sprintf("%s:disk:%s:%s", s.Config.Name, s.timestamp(), encodedPath), nil
	case model.Log:
		encodedPath := base64.StdEncoding.EncodeToString([]byte(value.Path))

		return fmt.Sprintf("%s:log:%s:%s", s.Config.Name, s.timestamp(), encodedPath), nil
	}

	return "", ErrUnknownModelType
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

func (s *Storage) alarmKey(m any) (string, error) {
	switch value := m.(type) {
	case model.CPU:
		return fmt.Sprintf("%s:alarm:cpu", s.Config.Name), nil
	case model.Memory:
		return fmt.Sprintf("%s:alarm:memory", s.Config.Name), nil
	case model.Load:
		return fmt.Sprintf("%s:alarm:load", s.Config.Name), nil
	case model.Disk:
		encodedPath := base64.StdEncoding.EncodeToString([]byte(value.Path))

		return fmt.Sprintf("%s:alarm:disk:%s", s.Config.Name, encodedPath), nil
	}

	return "", ErrUnknownModelType
}
