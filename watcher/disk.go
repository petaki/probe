package watcher

import (
	"log"

	"github.com/petaki/probe/model"
	"github.com/shirou/gopsutil/v4/disk"
)

func (w *Watcher) disk() {
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Println(err)

		return
	}

	var diskModels []model.Disk

	for _, partition := range partitions {
		if partition.Mountpoint == "" || w.storage.IsPathIgnored(partition.Mountpoint) {
			continue
		}

		diskUsage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			log.Println(err)

			continue
		}

		diskModels = append(diskModels, model.Disk{
			Path: partition.Mountpoint,
			Used: diskUsage.UsedPercent,
		})
	}

	for _, diskModel := range diskModels {
		err := w.storage.Save(diskModel)
		if err != nil {
			log.Println(err)
		}
	}
}
