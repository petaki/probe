package main // import "github.com/petaki/probe"

import (
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/petaki/probe/config"
	"github.com/petaki/probe/storage"
	"github.com/petaki/probe/watcher"
)

var (
	mainStorage = storage.Storage{}
)

func init() {
	fmt.Println("Starting Probe...")

	mainConfig, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	mainStorage = storage.New(&mainConfig)
}

func main() {
	fmt.Println("Probe is watching.")

	watcher.Watch(&mainStorage)
}
