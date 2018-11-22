package bootstrap

import (
	"fmt"
	"os"

	"github.com/petaki/probe/config"
)

// Config bootstrapper.
type Config struct{}

// Boot function.
func (Config) Boot() error {
	requiredKeys := config.GetRequiredKeys()

	for _, key := range requiredKeys {
		value, hasKey := os.LookupEnv(key)
		if !hasKey {
			return fmt.Errorf("%v is not defined", key)
		}

		err := config.Current.Parse(key, value)
		if err != nil {
			return err
		}
	}

	return nil
}
