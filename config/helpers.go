package config

import "strings"

func splitList(value string) []string {
	var entries []string

	for entry := range strings.SplitSeq(value, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		entries = append(entries, entry)
	}

	return entries
}
