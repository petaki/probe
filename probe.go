package main // import "github.com/petaki/probe"

import (
	"log"

	"github.com/petaki/probe/bootstrap"
	"github.com/petaki/probe/watchers"
)

func main() {
	err := bootstrap.Boot()
	if err != nil {
		log.Fatal(err)
	}

	err = watchers.Watch()
	if err != nil {
		log.Fatal(err)
	}
}
