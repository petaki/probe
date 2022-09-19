package watcher

import (
	"log"

	"github.com/petaki/probe/model"
	"github.com/petaki/probe/storage"
	"github.com/shirou/gopsutil/v3/disk"
)

// Disk watcher.
type Disk struct{}

// Watch function.
func (Disk) Watch(s *storage.Storage, index int, channel chan int) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Fatal(err)
	}

	var diskModels []model.Disk

	for _, partition := range partitions {
		diskUsage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			log.Fatal(err)
		}

		diskModels = append(diskModels, model.Disk{
			Path: partition.Mountpoint,
			Used: diskUsage.UsedPercent,
		})
	}

	for _, diskModel := range diskModels {
		err := s.Save(diskModel)
		if err != nil {
			log.Fatal(err)
		}
	}

	channel <- index
}
