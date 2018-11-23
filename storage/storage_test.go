package storage

import (
	"testing"

	"github.com/petaki/probe/config"
)

func TestNew(t *testing.T) {
	storage := New(config.Config{})

	if storage.Pool == nil {
		t.Errorf("The Pool is a nil pointer")
	}
}
