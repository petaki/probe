package watcher

import (
	"time"

	"github.com/petaki/probe/model"
	"github.com/shirou/gopsutil/cpu"
)

// CPU watcher.
type CPU struct{}

// Watch function.
func (CPU) Watch() (model.CPU, error) {
	cpuPercent, err := cpu.Percent(3*time.Second, false)
	if err != nil {
		return model.CPU{}, err
	}

	return model.CPU{
		Usage: cpuPercent[0],
	}, nil
}
