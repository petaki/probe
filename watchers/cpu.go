package watchers

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

// CPU watcher.
type CPU struct{}

// Watch function.
func (CPU) Watch() error {
	cpuPercent, err := cpu.Percent(3*time.Second, false)
	if err != nil {
		return err
	}

	fmt.Println(cpuPercent[0])

	return nil
}
