package svc

import (
	"github.com/goinbox/golog"
	"github.com/goinbox/gomisc"
)

type BaseSvc struct {
	accessLogger golog.ILogger
}

func NewBaseSvc() *BaseSvc {
	return &BaseSvc{
		accessLogger: new(golog.NoopLogger),
	}
}

func (b *BaseSvc) SetAccessLogger(logger golog.ILogger) *BaseSvc {
	if logger != nil {
		b.accessLogger = logger
	}

	return b
}

func (b *BaseSvc) AccessLogger() golog.ILogger {
	return b.accessLogger
}

func (b *BaseSvc) DebugLog(kind, msg []byte) {
	b.accessLogger.Debug(b.makeLogMsg(kind, msg))
}

func (b *BaseSvc) InfoLog(kind, msg []byte) {
	b.accessLogger.Info(b.makeLogMsg(kind, msg))
}

func (b *BaseSvc) NoticeLog(kind, msg []byte) {
	b.accessLogger.Notice(b.makeLogMsg(kind, msg))
}

func (b *BaseSvc) WarningLog(kind, msg []byte) {
	b.accessLogger.Warning(b.makeLogMsg(kind, msg))
}

func (b *BaseSvc) ErrorLog(kind, msg []byte) {
	b.accessLogger.Error(b.makeLogMsg(kind, msg))
}

func (b *BaseSvc) CriticalLog(kind, msg []byte) {
	b.accessLogger.Critical(b.makeLogMsg(kind, msg))
}

func (b *BaseSvc) AlertLog(kind, msg []byte) {
	b.accessLogger.Alert(b.makeLogMsg(kind, msg))
}

func (b *BaseSvc) EmergencyLog(kind, msg []byte) {
	b.accessLogger.Emergency(b.makeLogMsg(kind, msg))
}

func (b *BaseSvc) makeLogMsg(kind, msg []byte) []byte {
	return gomisc.AppendBytes(
		[]byte("["),
		kind,
		[]byte("]"),
		[]byte("\t"),
		msg,
	)
}
