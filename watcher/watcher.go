package watcher

import (
	"time"

	"github.com/petaki/probe/storage"
)

// Watcher interface.
type Watcher interface {
	Watch(s *storage.Storage, index int, channel chan int)
}

var watchers = []Watcher{
	CPU{},
	Disk{},
	Memory{},
	Process{},
}

// Watch function.
func Watch(s *storage.Storage) {
	channel := make(chan int)

	for i, watcher := range watchers {
		go watcher.Watch(s, i, channel)
	}

	for i := range channel {
		go func(index int) {
			time.Sleep(60 * time.Second)
			watchers[index].Watch(s, index, channel)
		}(i)
	}
}
