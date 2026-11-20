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

func startCPUTracking(p *process.Process) (float64, error) {
	_, err := p.Percent(0)
	if err != nil {
		return 0, err
	}

	return p.CPUPercent()
}
