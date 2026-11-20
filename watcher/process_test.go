package watcher

import (
	"os"
	"testing"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

func TestStartCPUTrackingSetsTheBaselineForTheNextSample(t *testing.T) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		t.Fatal(err)
	}

	_, err = startCPUTracking(p)
	if err != nil {
		t.Fatal(err)
	}

	deadline := time.Now().Add(300 * time.Millisecond)
	for x := 0; time.Now().Before(deadline); x++ {
		_ = x * x
	}

	used, err := p.Percent(0)
	if err != nil {
		t.Fatal(err)
	}

	if used <= 0 {
		t.Errorf("Expected the first sample after tracking starts to measure CPU use, got %.2f", used)
	}
}
