package svc

import (
	"github.com/goinbox/golog"
	"github.com/goinbox/gomisc"
)

type BaseSvc struct {
	AccessLogger golog.ILogger
}

func NewBaseSvc() *BaseSvc {
	return &BaseSvc{
		AccessLogger: new(golog.NoopLogger),
	}
}

func (b *BaseSvc) SetAccessLogger(logger golog.ILogger) *BaseSvc {
	if logger != nil {
		b.AccessLogger = logger
	}

	return b
}

func (b *BaseSvc) Debug(kind, msg []byte) {
	b.AccessLogger.Debug(b.makeLogMsg(kind, msg))
}

func (b *BaseSvc) Info(kind, msg []byte) {
	b.AccessLogger.Info(b.makeLogMsg(kind, msg))
}

func (b *BaseSvc) Notice(kind, msg []byte) {
	b.AccessLogger.Notice(b.makeLogMsg(kind, msg))
}

func (b *BaseSvc) Warning(kind, msg []byte) {
	b.AccessLogger.Warning(b.makeLogMsg(kind, msg))
}

func (b *BaseSvc) Error(kind, msg []byte) {
	b.AccessLogger.Error(b.makeLogMsg(kind, msg))
}

func (b *BaseSvc) Critical(kind, msg []byte) {
	b.AccessLogger.Critical(b.makeLogMsg(kind, msg))
}

func (b *BaseSvc) Alert(kind, msg []byte) {
	b.AccessLogger.Alert(b.makeLogMsg(kind, msg))
}

func (b *BaseSvc) Emergency(kind, msg []byte) {
	b.AccessLogger.Emergency(b.makeLogMsg(kind, msg))
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
