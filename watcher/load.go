package watcher

import (
	"log"

	"github.com/petaki/probe/model"
	"github.com/shirou/gopsutil/v4/load"
)

func (w *Watcher) load() {
	stat, err := load.Avg()
	if err != nil {
		log.Println(err)

		return
	}

	loadModel := model.Load{
		Load1:  stat.Load1,
		Load5:  stat.Load5,
		Load15: stat.Load15,
	}

	err = w.storage.Save(loadModel)
	if err != nil {
		log.Println(err)
	}
}
