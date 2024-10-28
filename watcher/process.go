package watcher

import (
	"cmp"
	"log"
	"slices"

	"github.com/petaki/probe/model"
	"github.com/petaki/probe/storage"
	"github.com/shirou/gopsutil/v3/process"
)

// Process watcher.
type Process struct{}

// Watch function.
func (Process) Watch(s *storage.Storage, index int, channel chan int) {
	processes, err := process.Processes()
	if err != nil {
		log.Fatal(err)
	}

	var processModels []model.Process

	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			name = "Unknown"
		}

		usedCPU, err := p.CPUPercent()
		if err != nil {
			usedCPU = 0
		}

		usedMemory, err := p.MemoryPercent()
		if err != nil {
			usedMemory = 0
		}

		processModels = append(processModels, model.Process{
			PID:        p.Pid,
			Name:       name,
			UsedCPU:    usedCPU,
			UsedMemory: usedMemory,
		})
	}

	slices.SortStableFunc(processModels, func(a, b model.Process) int {
		return cmp.Compare(b.UsedCPU, a.UsedCPU)
	})

	processModels = processModels[:3]

	err = s.Save(processModels)
	if err != nil {
		log.Fatal(err)
	}

	channel <- index
}
