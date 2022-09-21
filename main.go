package main

import (
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/petaki/probe/config"
	"github.com/petaki/probe/storage"
	"github.com/petaki/probe/watcher"
)

var (
	mainStorage = &storage.Storage{}
)

func init() {
	fmt.Println()
	fmt.Println("  🔍 Starting Probe...")
	fmt.Println()

	mainConfig, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	mainStorage = storage.New(mainConfig)

	if mainConfig.DataLogEnabled {
		fmt.Println("  📡 Data logging is enabled.")
		fmt.Println()
	}

	if mainConfig.AlarmEnabled {
		fmt.Println("  🚨 Alarm is armed.")
		fmt.Println()

		err = mainStorage.SaveAlarmConfig()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = mainStorage.DeleteAlarmConfig()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	fmt.Println("  🤖 Probe is watching.")
	fmt.Println()

	watcher.Watch(mainStorage)
}
