package storage

import (
	"testing"

	"github.com/petaki/probe/config"
)

func TestNew(t *testing.T) {
	storage := New(&config.Config{})

	if storage.Config == nil {
		t.Errorf("The config is a nil pointer")
	}

	if storage.Pool == nil {
		t.Errorf("The pool is a nil pointer")
	}
}
