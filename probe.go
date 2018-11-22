package main // import "github.com/petaki/probe"

import (
	"fmt"
	"log"

	"github.com/petaki/probe/bootstrap"
	"github.com/petaki/probe/config"
	"github.com/petaki/probe/watchers"
)

func main() {
	err := bootstrap.Boot()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(config.Current)

	err = watchers.Watch()
	if err != nil {
		log.Fatal(err)
	}
}
