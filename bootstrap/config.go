package bootstrap

import (
	"fmt"
	"os"
)

// Config bootstrapper.
type Config struct{}

// Boot function.
func (Config) Boot() error {
	requiredKeys := []string{
		"PROBE_PREFIX",
		"PROBE_KEEP_DATA",
		"PROBE_REDIS_HOST",
		"PROBE_REDIS_PASSWORD",
		"PROBE_REDIS_PORT",
		"PROBE_REDIS_DB",
	}

	for _, key := range requiredKeys {
		_, hasKey := os.LookupEnv(key)
		if !hasKey {
			return fmt.Errorf("%v is not defined", key)
		}
	}

	return nil
}
