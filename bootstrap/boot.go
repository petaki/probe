package bootstrap

// Bootstrapper interface.
type Bootstrapper interface {
	Boot() error
}

// Boot function.
func Boot() error {
	bootstrappers := []Bootstrapper{
		Print{},
		Dotenv{},
		Config{},
		Storage{},
	}

	for _, bootstrapper := range bootstrappers {
		err := bootstrapper.Boot()
		if err != nil {
			return err
		}
	}

	return nil
}
