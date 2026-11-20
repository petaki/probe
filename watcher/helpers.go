package watcher

import (
	"os"
	"strings"
)

func readTail(file *os.File, start, end int64) (string, error) {
	content := make([]byte, end-start)

	_, err := file.ReadAt(content, start)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(string(content), "\n"), nil
}
