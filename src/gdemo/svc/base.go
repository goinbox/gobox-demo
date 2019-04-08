package svc

import (
	"gdemo/misc"
	"gdemo/resource"
)

type BaseSvc struct {
	TraceId []byte
}

func (b *BaseSvc) DebugLog(point, msg []byte) {
	resource.AccessLogger.Debug(b.makeLogMsg(point, msg))
}

func (b *BaseSvc) InfoLog(point, msg []byte) {
	resource.AccessLogger.Info(b.makeLogMsg(point, msg))
}

func (b *BaseSvc) NoticeLog(point, msg []byte) {
	resource.AccessLogger.Notice(b.makeLogMsg(point, msg))
}

func (b *BaseSvc) WarningLog(point, msg []byte) {
	resource.AccessLogger.Warning(b.makeLogMsg(point, msg))
}

func (b *BaseSvc) ErrorLog(point, msg []byte) {
	resource.AccessLogger.Error(b.makeLogMsg(point, msg))
}

func (b *BaseSvc) CriticalLog(point, msg []byte) {
	resource.AccessLogger.Critical(b.makeLogMsg(point, msg))
}

func (b *BaseSvc) AlertLog(point, msg []byte) {
	resource.AccessLogger.Alert(b.makeLogMsg(point, msg))
}

func (b *BaseSvc) EmergencyLog(point, msg []byte) {
	resource.AccessLogger.Emergency(b.makeLogMsg(point, msg))
}

func (b *BaseSvc) makeLogMsg(point, msg []byte) []byte {
	return misc.FormatAccessLog(b.TraceId, point, msg)
}
