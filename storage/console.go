package storage

import (
	"fmt"

	"github.com/petaki/probe/model"
)

func (s *Storage) printValue(m any) error {
	switch value := m.(type) {
	case model.CPU:
		fmt.Printf("  ⚡ CPU: %.2f%%\n", value.Used)
		fmt.Println()
	case model.Memory:
		fmt.Printf("  📦 Memory: %.2f%%\n", value.Used)
		fmt.Println()
	case []model.ProcessCPU:
		for _, p := range value {
			fmt.Printf("  🚀 Process By CPU:[%d]%s: %.2f%%\n", p.PID, p.Name, p.Used)
			fmt.Println()
		}
	case []model.ProcessMemory:
		for _, p := range value {
			fmt.Printf("  🚀 Process By Memory:[%d]%s: %.2f%%\n", p.PID, p.Name, p.Used)
			fmt.Println()
		}
	case model.Load:
		fmt.Printf("  ✨ Load1: %.2f Load5: %.2f Load15: %.2f\n", value.Load1, value.Load5, value.Load15)
		fmt.Println()
	case model.Disk:
		fmt.Printf("  💾 Disk:%s: %.2f%%\n", value.Path, value.Used)
		fmt.Println()
	case model.Log:
		fmt.Printf("  📄 Log:%s:\n%s\n", value.Path, value.Content)
		fmt.Println()
	default:
		return ErrUnknownModelType
	}

	return nil
}
