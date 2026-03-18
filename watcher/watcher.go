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
	Memory{},
	Process{},
	Load{},
	Disk{},
	Log{},
}

// Watch function.
func Watch(s *storage.Storage) {
	channel := make(chan int)

	for i, watcher := range watchers {
		go watcher.Watch(s, i, channel)
	}

	for i := range channel {
		go func() {
			time.Sleep(60 * time.Second)
			watchers[i].Watch(s, i, channel)
		}()
	}
}
