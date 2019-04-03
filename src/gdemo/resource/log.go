package resource

import (
	"gdemo/conf"
	"gdemo/errno"

	"github.com/goinbox/exception"
	"github.com/goinbox/golog"
)

var AccessLogWriter golog.IWriter

var NoopLogger golog.ILogger = new(golog.NoopLogger)

var TestLogger golog.ILogger = golog.NewSimpleLogger(
	golog.NewConsoleWriter(),
	golog.NewConsoleFormater(golog.NewSimpleFormater())).
	SetLogLevel(golog.LEVEL_DEBUG)

func InitLog(systemName string) *exception.Exception {
	if conf.BaseConf.IsDev {
		AccessLogWriter = golog.NewConsoleWriter()
	} else {
		fw, err := golog.NewFileWriter(conf.LogConf.RootPath+"/"+systemName+"_access.log", conf.LogConf.Bufsize)
		if err != nil {
			return exception.New(errno.E_SYS_INIT_LOG_FAIL, err.Error())
		}
		AccessLogWriter = golog.NewAsyncWriter(fw, conf.LogConf.AsyncQueueSize)
	}

	return nil
}

func NewLogger(writer golog.IWriter, formater golog.IFormater) golog.ILogger {
	return golog.NewSimpleLogger(writer, formater).SetLogLevel(conf.LogConf.Level)
}

func FreeLog() {
	AccessLogWriter.Free()
}
