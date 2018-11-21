package watchers

// Watcher interface.
type Watcher interface {
	Watch() error
}

// Watch function.
func Watch() error {
	watchers := []Watcher{
		CPU{},
		Memory{},
		Disk{},
	}

	for _, watcher := range watchers {
		err := watcher.Watch()
		if err != nil {
			return err
		}
	}

	return nil
}
