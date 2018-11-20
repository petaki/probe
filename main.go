package main // import "github.com/petaki/probe"

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting Probe.")

	redisPrefix := os.Getenv("REDIS_PREFIX")
	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisPort := os.Getenv("REDIS_PORT")
	redisDb := os.Getenv("REDIS_DB")

	fmt.Println(redisPrefix, redisHost, redisPassword, redisPort, redisDb)

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
