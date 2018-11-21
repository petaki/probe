package main // import "github.com/petaki/probe"

import (
	"fmt"
	"log"
	"time"

	"github.com/petaki/probe/bootstrap"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	boot()

	cpuPercent, err := cpu.Percent(3*time.Second, false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cpuPercent[0])

	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(virtualMemory.UsedPercent)

	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Fatal(err)
	}

	for _, parititon := range partitions {
		usage, err := disk.Usage(parititon.Mountpoint)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(parititon.Mountpoint, usage.UsedPercent)
	}
}

func boot() {
	bootstrap.Dotenv()
	bootstrap.Config()
	bootstrap.Print()
}
