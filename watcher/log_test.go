package watcher

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/petaki/probe/config"
	"github.com/petaki/probe/storage"
)

func newTailWatcher(lines, bufferSize int) *Watcher {
	return New(&storage.Storage{
		Config: &config.Config{
			LogTailLines:      lines,
			LogTailBufferSize: bufferSize,
		},
	})
}

func writeTailFile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "app.log")

	err := os.WriteFile(path, []byte(content), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	return path
}

func TestTailFileReturnsLastLines(t *testing.T) {
	cases := []struct {
		name    string
		content string
		lines   int
		want    string
	}{
		{"empty file", "", 3, ""},
		{"fewer lines than requested", "a\nb\n", 3, "a\nb"},
		{"exactly the requested lines", "a\nb\nc\n", 3, "a\nb\nc"},
		{"more lines than requested", "a\nb\nc\nd\ne\n", 2, "d\ne"},
		{"no trailing newline", "a\nb\nc\nd", 2, "c\nd"},
		{"single line", "only line\n", 1, "only line"},
		{"single line without newline", "only line", 1, "only line"},
		{"last line only", "first\nsecond\nthird\n", 1, "third"},
		{"lines longer than the buffer", "alpha\nbravo\ncharlie\ndelta\n", 2, "charlie\ndelta"},
		{"blank lines in the middle", "a\n\n\nb\n", 3, "\n\nb"},
		{"multibyte characters", "árvíztűrő\ntükörfúrógép\n", 1, "tükörfúrógép"},
	}

	for _, tc := range cases {
		for _, bufferSize := range []int{1, 2, 3, 7, 4096} {
			t.Run(fmt.Sprintf("%s/buffer %d", tc.name, bufferSize), func(t *testing.T) {
				path := writeTailFile(t, tc.content)

				got, err := newTailWatcher(tc.lines, bufferSize).tailFile(path)
				if err != nil {
					t.Fatalf("tailFile returned error: %v", err)
				}

				if got != tc.want {
					t.Errorf("Expected %q, got %q", tc.want, got)
				}
			})
		}
	}
}

func TestTailFileCapsReadAtMaxSize(t *testing.T) {
	path := writeTailFile(t, "head\n"+strings.Repeat("x", tailMaxSize+10))

	got, err := newTailWatcher(10, 4096).tailFile(path)
	if err != nil {
		t.Fatalf("tailFile returned error: %v", err)
	}

	if len(got) != tailMaxSize {
		t.Errorf("Expected %d bytes, got %d", tailMaxSize, len(got))
	}

	if strings.Contains(got, "head") {
		t.Errorf("Expected content before the last %d bytes to be skipped", tailMaxSize)
	}
}

func TestTailFileReportsMissingFile(t *testing.T) {
	_, err := newTailWatcher(10, 4096).tailFile(filepath.Join(t.TempDir(), "missing.log"))
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("Expected fs.ErrNotExist, got %v", err)
	}
}
