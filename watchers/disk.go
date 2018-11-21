package watchers

import (
	"fmt"

	"github.com/shirou/gopsutil/disk"
)

// Disk watcher.
type Disk struct{}

// Watch function.
func (Disk) Watch() error {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return err
	}

	usage := map[string]float64{}

	for _, parititon := range partitions {
		diskUsage, err := disk.Usage(parititon.Mountpoint)
		if err != nil {
			return err
		}

		usage[parititon.Mountpoint] = diskUsage.UsedPercent
	}

	fmt.Println(usage)

	return nil
}
