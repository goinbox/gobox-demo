package golog

import (
	"os"
	"sync"
)

type ConsoleWriter struct {
	lock *sync.Mutex
}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{
		lock: new(sync.Mutex),
	}
}

func (c *ConsoleWriter) Write(msg []byte) (int, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	return os.Stdout.Write(msg)
}

func (c *ConsoleWriter) Flush() error {
	return nil
}

func (c *ConsoleWriter) Free() {
}
