package watcher

import (
	"log"
	"os"
	"strings"

	"github.com/petaki/probe/config"
	"github.com/petaki/probe/model"
	"github.com/petaki/probe/storage"
)

// Log watcher.
type Log struct{}

const tailMaxSize = 1024 * 1024

// Watch function.
func (Log) Watch(s *storage.Storage, index int, channel chan int) {
	defer func() { channel <- index }()

	if !s.Config.LogTailEnabled {
		return
	}

	for _, filePath := range s.Config.LogTailFiles {
		content, err := tailFile(filePath, s.Config)
		if err != nil {
			log.Println(err)

			continue
		}

		logModel := model.Log{
			Path:    filePath,
			Content: content,
		}

		err = s.Save(logModel)
		if err != nil {
			log.Println(err)
		}
	}
}

func tailFile(path string, c *config.Config) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", err
	}

	size := stat.Size()
	if size == 0 {
		return "", nil
	}

	limit := int64(0)
	if size > tailMaxSize {
		limit = size - tailMaxSize
	}

	newlines := 0
	cursor := size
	buf := make([]byte, c.LogTailBufferSize)

	for cursor > limit {
		chunkSize := min(int64(len(buf)), cursor-limit)
		cursor -= chunkSize

		n, err := file.ReadAt(buf[:chunkSize], cursor)
		if err != nil {
			return "", err
		}

		for i := n - 1; i >= 0; i-- {
			if buf[i] != '\n' {
				continue
			}

			position := cursor + int64(i) + 1

			if position == size && newlines == 0 {
				continue
			}

			newlines++

			if newlines >= c.LogTailLines {
				return readTail(file, position, size)
			}
		}
	}

	return readTail(file, limit, size)
}

func readTail(file *os.File, start, end int64) (string, error) {
	content := make([]byte, end-start)

	_, err := file.ReadAt(content, start)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(string(content), "\n"), nil
}
