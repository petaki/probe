package config

import "errors"

var (
	// ErrNotDefined error.
	ErrNotDefined = errors.New("config: not defined")

	// ErrInvalidKey error.
	ErrInvalidKey = errors.New("config: invalid key")

	// ErrInvalidTimeout error.
	ErrInvalidTimeout = errors.New("config: invalid timeout")

	// ErrInvalidPercent error.
	ErrInvalidPercent = errors.New("config: invalid percent")

	// ErrInvalidValue error.
	ErrInvalidValue = errors.New("config: invalid value")
)
