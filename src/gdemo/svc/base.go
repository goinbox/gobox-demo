package svc

import (
	"gdemo/misc"
	"gdemo/resource"
)

type BaseSvc struct {
	TraceId []byte
}

func (b *BaseSvc) DebugLog(point, msg []byte) {
	resource.AccessLogger.Debug(misc.FormatAccessLog(b.TraceId, point, msg))
}

func (b *BaseSvc) InfoLog(point, msg []byte) {
	resource.AccessLogger.Info(misc.FormatAccessLog(b.TraceId, point, msg))
}

func (b *BaseSvc) NoticeLog(point, msg []byte) {
	resource.AccessLogger.Notice(misc.FormatAccessLog(b.TraceId, point, msg))
}

func (b *BaseSvc) WarningLog(point, msg []byte) {
	resource.AccessLogger.Warning(misc.FormatAccessLog(b.TraceId, point, msg))
}

func (b *BaseSvc) ErrorLog(point, msg []byte) {
	resource.AccessLogger.Error(misc.FormatAccessLog(b.TraceId, point, msg))
}

func (b *BaseSvc) CriticalLog(point, msg []byte) {
	resource.AccessLogger.Critical(misc.FormatAccessLog(b.TraceId, point, msg))
}

func (b *BaseSvc) AlertLog(point, msg []byte) {
	resource.AccessLogger.Alert(misc.FormatAccessLog(b.TraceId, point, msg))
}

func (b *BaseSvc) EmergencyLog(point, msg []byte) {
	resource.AccessLogger.Emergency(misc.FormatAccessLog(b.TraceId, point, msg))
}
