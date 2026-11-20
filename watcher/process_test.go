package watcher

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/petaki/probe/config"
	"github.com/petaki/probe/storage"
	"github.com/shirou/gopsutil/v4/process"
)

func newProcessWatcher() *Watcher {
	return New(&storage.Storage{Config: &config.Config{}})
}

func burnCPU(d time.Duration) {
	deadline := time.Now().Add(d)
	for x := 0; time.Now().Before(deadline); x++ {
		_ = x * x
	}
}

func startBusyChild(t *testing.T) int32 {
	t.Helper()

	child := exec.Command("sh", "-c", "while :; do :; done")

	err := child.Start()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = child.Process.Kill()
		_ = child.Wait()
	})

	return int32(child.Process.Pid)
}

func TestTrackCPUMeasuresTheNextSampleOfANewProcess(t *testing.T) {
	w := newProcessWatcher()
	pid := int32(os.Getpid())

	first, err := process.NewProcess(pid)
	if err != nil {
		t.Fatal(err)
	}

	tracked, _, err := w.trackCPU(first)
	if err != nil {
		t.Fatal(err)
	}

	w.processes[pid] = tracked

	burnCPU(300 * time.Millisecond)

	next, err := process.NewProcess(pid)
	if err != nil {
		t.Fatal(err)
	}

	_, used, err := w.trackCPU(next)
	if err != nil {
		t.Fatal(err)
	}

	if used <= 0 {
		t.Errorf("Expected the second sample to measure CPU use, got %.2f", used)
	}
}

func TestTrackCPUStartsOverWhenAPIDIsReused(t *testing.T) {
	w := newProcessWatcher()

	stale, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		t.Fatal(err)
	}

	burnCPU(1500 * time.Millisecond)

	_, err = stale.Percent(0)
	if err != nil {
		t.Fatal(err)
	}

	pid := startBusyChild(t)

	time.Sleep(time.Second)

	stale.Pid = pid
	w.processes[pid] = stale

	fresh, err := process.NewProcess(pid)
	if err != nil {
		t.Fatal(err)
	}

	tracked, used, err := w.trackCPU(fresh)
	if err != nil {
		t.Fatal(err)
	}

	if used <= 0 {
		t.Errorf("Expected the busy process to report CPU use, got %.2f", used)
	}

	if tracked != fresh {
		t.Errorf("Expected the reused PID to be tracked as a new process")
	}
}
