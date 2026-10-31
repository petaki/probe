package watcher

import "errors"

var (
	// ErrNoCPUPercent error.
	ErrNoCPUPercent = errors.New("watcher: no cpu percent")
)
