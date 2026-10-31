package watcher

import (
	"time"

	"github.com/petaki/probe/storage"
	"github.com/shirou/gopsutil/v4/process"
)

// Watcher type.
type Watcher struct {
	storage   *storage.Storage
	processes map[int32]*process.Process
}

// New function.
func New(s *storage.Storage) *Watcher {
	return &Watcher{
		storage:   s,
		processes: map[int32]*process.Process{},
	}
}

// Watch function.
func (w *Watcher) Watch() {
	time.Sleep(time.Second)

	for _, watch := range []func(){w.cpu, w.memory, w.process, w.load, w.disk, w.log} {
		go func() {
			for {
				watch()
				time.Sleep(time.Until(time.Now().Truncate(time.Minute).Add(time.Minute)))
			}
		}()
	}

	select {}
}
