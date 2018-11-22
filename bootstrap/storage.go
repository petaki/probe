package bootstrap

import (
	"github.com/petaki/probe/storage"
)

// Storage bootstrapper.
type Storage struct{}

// Boot function.
func (Storage) Boot() error {
	storage.Current.Setup()

	return nil
}
