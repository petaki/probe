package watcher

import (
	"cmp"
	"log"
	"slices"

	"github.com/petaki/probe/model"
	"github.com/petaki/probe/storage"
	"github.com/shirou/gopsutil/v4/process"
)

// Process watcher.
type Process struct{}

const processTopN = 3

// Watch function.
func (Process) Watch(s *storage.Storage, index int, channel chan int) {
	defer func() { channel <- index }()

	processes, err := process.Processes()
	if err != nil {
		log.Println(err)

		return
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

	if len(processCPUModels) > processTopN {
		processCPUModels = processCPUModels[:processTopN]
	}

	if len(processMemoryModels) > processTopN {
		processMemoryModels = processMemoryModels[:processTopN]
	}

	err = s.Save(processCPUModels)
	if err != nil {
		log.Println(err)
	}

	err = s.Save(processMemoryModels)
	if err != nil {
		log.Println(err)
	}
}
