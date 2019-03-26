/**
* @file file.go
* @brief writer msg to file
* @author ligang
* @date 2016-02-03
 */

package golog

import (
	"github.com/goinbox/gomisc"

	"os"
	"sync"
	"time"
)

type FileWriter struct {
	*os.File

	path           string
	lock           *sync.Mutex
	lastTimeSecond int64

	buf     []byte
	bufsize int
	bufpos  int
}

func NewFileWriter(path string, bufsize int) (*FileWriter, error) {
	file, err := openFile(path)
	if err != nil {
		return nil, err
	}

	return &FileWriter{
		File: file,

		path:           path,
		lock:           new(sync.Mutex),
		lastTimeSecond: time.Now().Unix(),

		buf:     make([]byte, bufsize),
		bufsize: bufsize,
		bufpos:  0,
	}, nil
}

func (f *FileWriter) Write(msg []byte) (int, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	err := f.ensureFileExist()
	if err != nil {
		return 0, err
	}

	if f.bufsize == 0 {
		return f.File.Write(msg)
	}

	if f.appendToBuffer(msg) {
		return len(msg), nil
	}

	err = f.flushBuffer()
	if err != nil {
		return 0, err
	}

	if f.appendToBuffer(msg) {
		return len(msg), nil
	}

	return f.File.Write(msg)
}

func (f *FileWriter) Flush() error {
	if f.bufsize == 0 {
		return nil
	}

	f.lock.Lock()
	defer f.lock.Unlock()

	err := f.ensureFileExist()
	if err != nil {
		return err
	}

	return f.flushBuffer()
}

func (f *FileWriter) Free() {
	f.ensureFileExist()

	f.flushBuffer()
	f.File.Close()
}

func openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
}

func (f *FileWriter) ensureFileExist() error {
	nowTimeSecond := time.Now().Unix()
	if f.lastTimeSecond == nowTimeSecond {
		return nil
	}

	if gomisc.FileExist(f.path) {
		return nil
	}

	f.Close()

	var err error
	f.File, err = openFile(f.path)
	if err != nil {
		return err
	}

	f.lastTimeSecond = nowTimeSecond
	return nil
}

func (f *FileWriter) appendToBuffer(msg []byte) bool {
	after := f.bufpos + len(msg)
	if after >= f.bufsize {
		return false
	}

	copy(f.buf[f.bufpos:], msg)
	f.bufpos = after

	return true
}

func (f *FileWriter) flushBuffer() error {
	if f.bufpos == 0 {
		return nil
	}

	_, err := f.File.Write(f.buf[:f.bufpos])
	f.bufpos = 0

	return err
}
