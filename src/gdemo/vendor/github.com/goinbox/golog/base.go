/**
* @file logger.go
* @author ligang
* @date 2016-02-04
 */

package golog

import "io"

const (
	LEVEL_EMERGENCY = 0
	LEVEL_ALERT     = 1
	LEVEL_CRITICAL  = 2
	LEVEL_ERROR     = 3
	LEVEL_WARNING   = 4
	LEVEL_NOTICE    = 5
	LEVEL_INFO      = 6
	LEVEL_DEBUG     = 7
)

var LogLevels map[int][]byte = map[int][]byte{
	LEVEL_EMERGENCY: []byte("emergency"),
	LEVEL_ALERT:     []byte("alert"),
	LEVEL_CRITICAL:  []byte("critical"),
	LEVEL_ERROR:     []byte("error"),
	LEVEL_WARNING:   []byte("warning"),
	LEVEL_NOTICE:    []byte("notice"),
	LEVEL_INFO:      []byte("info"),
	LEVEL_DEBUG:     []byte("debug"),
}

type ILogger interface {
	Debug(msg []byte)
	Info(msg []byte)
	Notice(msg []byte)
	Warning(msg []byte)
	Error(msg []byte)
	Critical(msg []byte)
	Alert(msg []byte)
	Emergency(msg []byte)

	Log(level int, msg []byte) error

	Flush() error
	Free()
}

type IFormater interface {
	Format(level int, msg []byte) []byte
}

type IWriter interface {
	io.Writer

	Flush() error
	Free()
}
