package watcher

import (
	"log"

	"github.com/petaki/probe/model"
	"github.com/petaki/probe/storage"
	"github.com/shirou/gopsutil/v3/cpu"
)

// CPU watcher.
type CPU struct{}

// Watch function.
func (CPU) Watch(s *storage.Storage, index int, channel chan int) {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		log.Fatal(err)
	}

	cpuModel := model.CPU{
		Used: cpuPercent[0],
	}

	err = s.Save(cpuModel)
	if err != nil {
		log.Fatal(err)
	}

	channel <- index
}
