package watcher

import (
	"log"

	"github.com/petaki/probe/model"
	"github.com/petaki/probe/storage"
	"github.com/shirou/gopsutil/v4/mem"
)

// Memory watcher.
type Memory struct{}

// Watch function.
func (Memory) Watch(s *storage.Storage, index int, channel chan int) {
	defer func() { channel <- index }()

	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		log.Println(err)

		return
	}

	memoryModel := model.Memory{
		Used: virtualMemory.UsedPercent,
	}

	err = s.Save(memoryModel)
	if err != nil {
		log.Println(err)
	}
}
