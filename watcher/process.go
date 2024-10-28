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

	var processCPUModels []model.ProcessCPU
	var processMemoryModels []model.ProcessMemory

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

		processCPUModels = append(processCPUModels, model.ProcessCPU{
			PID:  p.Pid,
			Name: name,
			Used: usedCPU,
		})

		processMemoryModels = append(processMemoryModels, model.ProcessMemory{
			PID:  p.Pid,
			Name: name,
			Used: usedMemory,
		})
	}

	slices.SortStableFunc(processCPUModels, func(a, b model.ProcessCPU) int {
		return cmp.Compare(b.Used, a.Used)
	})

	slices.SortStableFunc(processMemoryModels, func(a, b model.ProcessMemory) int {
		return cmp.Compare(b.Used, a.Used)
	})

	processCPUModels = processCPUModels[:3]
	processMemoryModels = processMemoryModels[:3]

	err = s.Save(processCPUModels)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Save(processMemoryModels)
	if err != nil {
		log.Fatal(err)
	}

	channel <- index
}
