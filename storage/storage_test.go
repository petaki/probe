package storage

import (
	"testing"

	"github.com/petaki/probe/config"
)

func TestNew(t *testing.T) {
	storage := New(&config.Config{
		RedisPrefix:  "probe:",
		RedisTimeout: 10080,
	})

	if storage.Prefix != "probe:" {
		t.Errorf("Expected prefix probe:, but got %v", storage.Prefix)
	}

	if storage.Timeout != 10080 {
		t.Errorf("Expected timeout 10080, but got %v", storage.Timeout)
	}

	if storage.Pool == nil {
		t.Errorf("The pool is a nil pointer")
	}
}
