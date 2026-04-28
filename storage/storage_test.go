package storage

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/petaki/probe/config"
	"github.com/petaki/probe/model"
)

func TestNew(t *testing.T) {
	storage := New(&config.Config{})

	if storage.Config == nil {
		t.Errorf("The config is a nil pointer")
	}

	if storage.Client != nil {
		t.Errorf("The client is not a nil pointer")
	}

	if storage.Pool != nil {
		t.Errorf("The pool is not a nil pointer")
	}
}

func TestIsPathIgnored(t *testing.T) {
	storage := &Storage{
		Config: &config.Config{
			DiskIgnored: []string{"/dev", "/var/lib/docker/*", "*tmpfs*", "*.snap"},
		},
	}

	cases := []struct {
		path string
		want bool
	}{
		{"/dev", true},
		{"/dev/sda", false},
		{"/var/lib/docker/overlay", true},
		{"/var/lib/docker", false},
		{"/run/tmpfs/x", true},
		{"/snap/core/foo.snap", true},
		{"/foo.snap", true},
		{"/", false},
		{"/home", false},
	}

	for _, tc := range cases {
		got := storage.isPathIgnored(tc.path)
		if got != tc.want {
			t.Errorf("isPathIgnored(%q) = %v, want %v", tc.path, got, tc.want)
		}
	}
}

func TestCallAlarmInterpolatesProbeName(t *testing.T) {
	var receivedBody string
	var receivedMethod string
	var receivedAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedAuth = r.Header.Get("Authorization")

		body, _ := io.ReadAll(r.Body)
		receivedBody = string(body)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	storage := &Storage{
		Config: &config.Config{
			Name:               "myhost",
			AlarmCPUPercent:    30,
			AlarmWebhookMethod: "POST",
			AlarmWebhookURL:    server.URL,
			AlarmWebhookHeader: map[string]string{"Authorization": "Bearer TOKEN"},
			AlarmWebhookBody:   `{"probe":"%p","name":"%n","alarm":%a,"used":%u,"link":"%l"}`,
		},
		Client: &http.Client{Timeout: 5 * time.Second},
	}

	err := storage.callAlarm(model.CPU{Used: 90.5})
	if err != nil {
		t.Fatalf("callAlarm returned error: %v", err)
	}

	if receivedMethod != "POST" {
		t.Errorf("Expected POST method, got %v", receivedMethod)
	}

	if receivedAuth != "Bearer TOKEN" {
		t.Errorf("Expected Authorization header to be set, got %v", receivedAuth)
	}

	if !strings.Contains(receivedBody, `"probe":"myhost"`) {
		t.Errorf("Expected body to contain probe=myhost, got %v", receivedBody)
	}

	if !strings.Contains(receivedBody, `"name":"CPU"`) {
		t.Errorf("Expected body to contain name=CPU, got %v", receivedBody)
	}

	if !strings.Contains(receivedBody, `"link":"/cpu?probe=myhost"`) {
		t.Errorf("Expected body to contain link with probe=myhost, got %v", receivedBody)
	}
}
