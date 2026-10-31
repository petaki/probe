package watcher

import (
	"cmp"
	"log"
	"slices"

	"github.com/petaki/probe/model"
	"github.com/shirou/gopsutil/v4/process"
)

const processTopN = 3

func (w *Watcher) process() {
	processes, err := process.Processes()
	if err != nil {
		log.Println(err)

		return
	}

	var processCPUModels []model.ProcessCPU
	var processMemoryModels []model.ProcessMemory

	current := make(map[int32]*process.Process, len(processes))

	for _, p := range processes {
		previous, seen := w.processes[p.Pid]
		if !seen {
			previous = p
		}

		current[p.Pid] = previous

		name, err := p.Name()
		if err != nil {
			name = "Unknown"
		}

		usedCPU, err := previous.Percent(0)
		if !seen {
			usedCPU, err = p.CPUPercent()
		}

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

	w.processes = current

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

	err = w.storage.Save(processCPUModels)
	if err != nil {
		log.Println(err)
	}

	err = w.storage.Save(processMemoryModels)
	if err != nil {
		log.Println(err)
	}
}
