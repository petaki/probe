package model

// Process model.
type Process struct {
	PID        int32
	Name       string
	UsedCPU    float64
	UsedMemory float32
}
