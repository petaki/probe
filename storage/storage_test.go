package storage

import (
	"testing"

	"github.com/petaki/probe/config"
)

func TestNew(t *testing.T) {
	storage := New(&config.Config{
		RedisPrefix: "probe:",
	})

	if storage.Prefix != "probe:" {
		t.Errorf("Expected prefix probe:, but got %v", storage.Prefix)
	}

	if storage.Pool == nil {
		t.Errorf("The pool is a nil pointer")
	}
}
