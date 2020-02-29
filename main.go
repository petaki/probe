package main // import "github.com/petaki/probe"

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/petaki/probe/config"
	"github.com/petaki/probe/storage"
	"github.com/petaki/probe/watcher"
)

var (
	mainConfig  = config.Config{}
	mainStorage = storage.Storage{}
)

func init() {
	fmt.Println("Starting Probe...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	mainConfig, err = config.Load()
	if err != nil {
		log.Fatal(err)
	}

	mainStorage = storage.New(&mainConfig)
}

func main() {
	fmt.Println("Probe is watching.")

	watcher.Watch(&mainStorage)
}
