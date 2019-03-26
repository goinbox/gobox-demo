package svc

import (
	"github.com/goinbox/golog"
	"github.com/goinbox/gomisc"
)

type BaseSvc struct {
	TraceId      []byte
	AccessLogger golog.ILogger
}

func (b *BaseSvc) DebugLog(point, msg []byte) {
	b.AccessLogger.Debug(b.makeLogMsg(point, msg))
}

func (b *BaseSvc) InfoLog(point, msg []byte) {
	b.AccessLogger.Info(b.makeLogMsg(point, msg))
}

func (b *BaseSvc) NoticeLog(point, msg []byte) {
	b.AccessLogger.Notice(b.makeLogMsg(point, msg))
}

func (b *BaseSvc) WarningLog(point, msg []byte) {
	b.AccessLogger.Warning(b.makeLogMsg(point, msg))
}

func (b *BaseSvc) ErrorLog(point, msg []byte) {
	b.AccessLogger.Error(b.makeLogMsg(point, msg))
}

func (b *BaseSvc) CriticalLog(point, msg []byte) {
	b.AccessLogger.Critical(b.makeLogMsg(point, msg))
}

func (b *BaseSvc) AlertLog(point, msg []byte) {
	b.AccessLogger.Alert(b.makeLogMsg(point, msg))
}

func (b *BaseSvc) EmergencyLog(point, msg []byte) {
	b.AccessLogger.Emergency(b.makeLogMsg(point, msg))
}

func (b *BaseSvc) makeLogMsg(point, msg []byte) []byte {
	return gomisc.AppendBytes(
		[]byte("["),
		point,
		[]byte("]"),
		[]byte("\t"),
		msg,
	)
}
