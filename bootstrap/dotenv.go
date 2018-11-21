package bootstrap

import (
	"github.com/joho/godotenv"
)

// Dotenv bootstrapper.
type Dotenv struct{}

// Boot function.
func (Dotenv) Boot() error {
	return godotenv.Load()
}
