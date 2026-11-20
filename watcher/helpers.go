package watcher

import (
	"os"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

func readTail(file *os.File, start, end int64) (string, error) {
	content := make([]byte, end-start)

	_, err := file.ReadAt(content, start)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(string(content), "\n"), nil
}

func sameStartTime(a, b *process.Process) bool {
	aCreated, err := a.CreateTime()
	if err != nil {
		return false
	}

	bCreated, err := b.CreateTime()
	if err != nil {
		return false
	}

	return aCreated == bCreated
}
