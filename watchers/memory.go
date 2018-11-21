package watchers

import (
	"fmt"

	"github.com/shirou/gopsutil/mem"
)

// Memory watcher.
type Memory struct{}

// Watch function.
func (Memory) Watch() error {
	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	fmt.Println(virtualMemory.UsedPercent)

	return nil
}
