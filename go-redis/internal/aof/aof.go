package aof

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

type AOF struct {
	file *os.File
	mu   sync.Mutex
}

func New(path string) (*AOF, error) {
	file, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_APPEND|os.O_RDWR,
		0644,
	)

	if err != nil {
		return nil, err
	}

	return &AOF{file: file}, nil
}

func (aof *AOF) Append(cmd []string) error {
	aof.mu.Lock()
	defer aof.mu.Lock()

	line := strings.Join(cmd, " ") + "\n"
	_, err := aof.file.WriteString(line)
	if err != nil {
		return err
	}

	return aof.file.Sync()
}

func (aof *AOF) Replay(apply func([]string)) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Seek(0, 0)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(aof.file)
	for scanner.Scan() {
		line := scanner.Text()
		cmd := strings.Split(line, " ")
		apply(cmd)
	}

	return scanner.Err()
}
