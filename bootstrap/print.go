package bootstrap

import "fmt"

// Print bootstrapper.
type Print struct{}

// Boot function.
func (Print) Boot() error {
	fmt.Println("Starting Probe.")

	return nil
}
