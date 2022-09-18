package watcher

import (
	"log"

	"github.com/petaki/probe/model"
	"github.com/petaki/probe/storage"
	"github.com/shirou/gopsutil/v3/mem"
)

// Memory watcher.
type Memory struct{}

// Watch function.
func (Memory) Watch(s *storage.Storage, index int, channel chan int) {
	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}

	memoryModel := model.Memory{
		Used: virtualMemory.UsedPercent,
	}

	err = s.Save(memoryModel)
	if err != nil {
		log.Fatal(err)
	}

	channel <- index
}
