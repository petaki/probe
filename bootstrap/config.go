package bootstrap

import (
	"log"
	"os"
)

// Config function.
func Config() {
	requiredKeys := []string{
		"PROBE_PREFIX",
		"PROBE_KEEP_DATA",
		"PROBE_REDIS_HOST",
		"PROBE_REDIS_PASSWORD",
		"PROBE_REDIS_PORT",
		"PROBE_REDIS_DB",
	}

	var hasKey bool

	for _, key := range requiredKeys {
		_, hasKey = os.LookupEnv(key)
		if !hasKey {
			log.Fatal(key, " is not defined")
		}
	}
}
