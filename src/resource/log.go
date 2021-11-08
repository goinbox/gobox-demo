package resource

import (
	"fmt"

	"github.com/goinbox/golog"

	"gdemo/conf"
)

var accessLogWriter golog.Writer
var AccessLogger golog.Logger

func InitLog(config *conf.LogConf) error {
	var err error

	accessLogWriter, err = golog.NewFileWriter(config.Path, config.Bufsize)
	if err != nil {
		return fmt.Errorf("golog.NewFileWriter error: %w", err)
	}

	if config.AsyncQueueSize > 0 {
		accessLogWriter = golog.NewAsyncWriter(accessLogWriter, config.AsyncQueueSize)
	}

	logger := golog.NewSimpleLogger(accessLogWriter, formater(config)).SetLogLevel(config.Level)
	if config.EnableColor {
		logger.EnableColor()
	}

	AccessLogger = logger

	return nil
}

func FreeLog() {
	accessLogWriter.Free()
}

func formater(config *conf.LogConf) golog.Formater {
	switch config.Formater {
	case "simple":
		return golog.NewSimpleFormater()
	}

	return golog.NewJsonFormater()
}
