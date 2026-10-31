package watcher

import (
	"log"
	"os"
	"strings"

	"github.com/petaki/probe/model"
)

const tailMaxSize = 1024 * 1024

func (w *Watcher) log() {
	if !w.storage.Config.LogTailEnabled {
		return
	}

	for _, filePath := range w.storage.Config.LogTailFiles {
		content, err := w.tailFile(filePath)
		if err != nil {
			log.Println(err)

			continue
		}

		logModel := model.Log{
			Path:    filePath,
			Content: content,
		}

		err = w.storage.Save(logModel)
		if err != nil {
			log.Println(err)
		}
	}
}

func (w *Watcher) tailFile(path string) (string, error) {
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
	buf := make([]byte, w.storage.Config.LogTailBufferSize)

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

			if newlines >= w.storage.Config.LogTailLines {
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
