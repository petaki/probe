package watcher

import (
	"log"

	"github.com/petaki/probe/model"
	"github.com/petaki/probe/storage"
	"github.com/shirou/gopsutil/v4/load"
)

// Load watcher.
type Load struct{}

// Watch function.
func (Load) Watch(s *storage.Storage, index int, channel chan int) {
	defer func() { channel <- index }()

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

	err = s.Save(loadModel)
	if err != nil {
		log.Println(err)
	}
}
