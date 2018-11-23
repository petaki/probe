package watcher

import (
	"github.com/petaki/probe/model"
	"github.com/shirou/gopsutil/disk"
)

// Disk watcher.
type Disk struct{}

// Watch function.
func (Disk) Watch() (model.Disk, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return model.Disk{}, err
	}

	diskModel := model.Disk{}

	for _, parititon := range partitions {
		diskUsage, err := disk.Usage(parititon.Mountpoint)
		if err != nil {
			return model.Disk{}, err
		}

		diskModel.Usage[parititon.Mountpoint] = diskUsage.UsedPercent
	}

	return diskModel, nil
}
