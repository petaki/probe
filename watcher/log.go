package watcher

import (
	"log"
	"os"
	"strings"

	"github.com/petaki/probe/model"
	"github.com/petaki/probe/storage"
)

// Log watcher.
type Log struct{}

// Watch function.
func (Log) Watch(s *storage.Storage, index int, channel chan int) {
	if !s.Config.LogTailEnabled {
		channel <- index
		return
	}

	for _, filePath := range s.Config.LogTailFiles {
		content, err := tailFile(filePath, s.Config.LogTailLines, s.Config.LogTailBufferSize)
		if err != nil {
			log.Fatal(err)
		}

		logModel := model.Log{
			Path:    filePath,
			Content: content,
		}

		err = s.Save(logModel)
		if err != nil {
			log.Fatal(err)
		}
	}

	channel <- index
}

func tailFile(path string, lines int, bufferSize int) (string, error) {
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

	newlines := 0
	cursor := size
	buf := make([]byte, bufferSize)

	for cursor > 0 {
		chunkSize := min(int64(len(buf)), cursor)
		cursor -= chunkSize

		n, err := file.ReadAt(buf[:chunkSize], cursor)
		if err != nil {
			return "", err
		}

		for i := n - 1; i >= 0; i-- {
			if buf[i] != '\n' {
				continue
			}

			if cursor+int64(i)+1 == size && newlines == 0 {
				continue
			}

			newlines++

			if newlines >= lines {
				start := cursor + int64(i) + 1
				content := make([]byte, size-start)

				_, err = file.ReadAt(content, start)
				if err != nil {
					return "", err
				}

				return strings.TrimRight(string(content), "\n"), nil
			}
		}
	}

	content := make([]byte, size)

	_, err = file.ReadAt(content, 0)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(string(content), "\n"), nil
}
