package watcher

import (
	"log"

	"github.com/petaki/probe/model"
	"github.com/shirou/gopsutil/v4/mem"
)

func (w *Watcher) memory() {
	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		log.Println(err)

		return
	}

	memoryModel := model.Memory{
		Used: virtualMemory.UsedPercent,
	}

	err = w.storage.Save(memoryModel)
	if err != nil {
		log.Println(err)
	}
}
