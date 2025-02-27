package config

import "errors"

var (
	// ErrInvalidTimeout error.
	ErrInvalidTimeout = errors.New("config: invalid timeout")

	// ErrInvalidPercent error.
	ErrInvalidPercent = errors.New("config: invalid percent")

	// ErrInvalidValue error.
	ErrInvalidValue = errors.New("config: invalid value")
)
