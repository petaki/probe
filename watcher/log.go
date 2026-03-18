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
		content, err := tailFile(filePath, s.Config.LogTailLines)
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

func tailFile(path string, lines int) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	allLines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")

	if len(allLines) <= lines {
		return strings.Join(allLines, "\n"), nil
	}

	return strings.Join(allLines[len(allLines)-lines:], "\n"), nil
}
