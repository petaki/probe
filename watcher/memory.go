package watcher

import (
	"github.com/petaki/probe/model"
	"github.com/shirou/gopsutil/mem"
)

// Memory watcher.
type Memory struct{}

// Watch function.
func (Memory) Watch() (model.Memory, error) {
	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		return model.Memory{}, err
	}

	return model.Memory{
		Usage: virtualMemory.UsedPercent,
	}, nil
}
