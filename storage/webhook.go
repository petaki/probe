package storage

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/petaki/probe/model"
)

func (s *Storage) callAlarm(m any) error {
	probe := s.Config.Name

	var name string
	var used string
	var alarm float64
	var link string

	switch value := m.(type) {
	case model.CPU:
		name = "CPU"
		alarm = s.Config.AlarmCPUPercent
		used = fmt.Sprintf("%.2f", value.Used)
		link = fmt.Sprintf("/cpu?probe=%s", probe)
	case model.Memory:
		name = "Memory"
		alarm = s.Config.AlarmMemoryPercent
		used = fmt.Sprintf("%.2f", value.Used)
		link = fmt.Sprintf("/memory?probe=%s", probe)
	case []model.ProcessCPU:
		return nil
	case []model.ProcessMemory:
		return nil
	case model.Load:
		name = "Load"
		alarm = s.Config.AlarmLoadValue
		used = fmt.Sprintf("\"%.2f,%.2f,%.2f\"", value.Load1, value.Load5, value.Load15)
		link = fmt.Sprintf("/load?probe=%s", probe)
	case model.Disk:
		name = fmt.Sprintf("Disk:%s", value.Path)
		alarm = s.Config.AlarmDiskPercent
		used = fmt.Sprintf("%.2f", value.Used)
		link = fmt.Sprintf("/disk?probe=%s&path=%s", probe, value.Path)
	case model.Log:
		return nil
	default:
		return ErrUnknownModelType
	}

	now := time.Now()

	body := strings.ReplaceAll(s.Config.AlarmWebhookBody, "%p", probe)
	body = strings.ReplaceAll(body, "%n", name)
	body = strings.ReplaceAll(body, "%a", fmt.Sprintf("%.2f", alarm))
	body = strings.ReplaceAll(body, "%u", used)
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
