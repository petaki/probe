package storage

import "errors"

var (
	// ErrBadStatusCode error.
	ErrBadStatusCode = errors.New("storage: bad status code")

	// ErrUnknownModelType error.
	ErrUnknownModelType = errors.New("storage: unknown model type")
)
