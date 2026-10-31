package watcher

import (
	"log"

	"github.com/petaki/probe/model"
	"github.com/shirou/gopsutil/v4/cpu"
)

func (w *Watcher) cpu() {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		log.Println(err)

		return
	}

	if len(cpuPercent) == 0 {
		log.Println(ErrNoCPUPercent)

		return
	}

	cpuModel := model.CPU{
		Used: cpuPercent[0],
	}

	err = w.storage.Save(cpuModel)
	if err != nil {
		log.Println(err)
	}
}
