package resource

import (
	"gdemo/conf"
	"gdemo/perror"

	"gdemo/perror"
	"github.com/goinbox/golog"
)

var accessLogWriter golog.IWriter
var AccessLogger golog.ILogger

var traceLogWriter golog.IWriter
var TraceLogger golog.ILogger

var NoopLogger golog.ILogger = new(golog.NoopLogger)

var TestLogger golog.ILogger = golog.NewSimpleLogger(
	golog.NewConsoleWriter(),
	golog.NewConsoleFormater(golog.NewSimpleFormater())).
	SetLogLevel(golog.LevelDebug)

func InitLog(systemName string) *perror.Error {
	if conf.BaseConf.IsDev {
		accessLogWriter = golog.NewConsoleWriter()
	} else {
		fw, err := golog.NewFileWriter(conf.LogConf.RootPath+"/"+systemName+"_access.log", conf.LogConf.Bufsize)
		if err != nil {
			return perror.Error(perror.ESysInitLogFail, err.Error())
		}
		accessLogWriter = golog.NewAsyncWriter(fw, conf.LogConf.AsyncQueueSize)
	}
	AccessLogger = NewLogger(accessLogWriter)

	fw, err := golog.NewFileWriter(conf.LogConf.RootPath+"/"+systemName+"_trace.log", conf.LogConf.Bufsize)
	if err != nil {
		return perror.Error(perror.ESysInitLogFail, err.Error())
	}
	traceLogWriter = golog.NewAsyncWriter(fw, conf.LogConf.AsyncQueueSize)
	TraceLogger = NewLogger(traceLogWriter)

	return nil
}

func NewLogger(writer golog.IWriter) golog.ILogger {
	return golog.NewSimpleLogger(writer, golog.NewSimpleFormater()).SetLogLevel(conf.LogConf.Level)
}

func FreeLog() {
	accessLogWriter.Free()
	traceLogWriter.Free()
}
