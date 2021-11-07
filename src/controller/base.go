package controller

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	gcontroller "github.com/goinbox/gohttp/controller"

	"gdemo/misc"
	"gdemo/resource"
)

const (
	RemoteRealIpHeaderKey   = "REMOTE-REAL-IP"
	RemoteRealPortHeaderKey = "REMOTE-REAL-PORT"

	DownstreamServerIp = "127.0.0.1"

	MaxAutoParseBodyLen = 1024 * 1024

	TraceIdHeaderKey = "TRACE-ID"
	TraceIdQueryKey  = "traceId"
)

type BaseContext struct {
	*gcontroller.BaseContext

	ReqRawBody []byte

	QueryValues    url.Values
	RemoteRealAddr struct {
		Ip   string
		Port string
	}

	TraceId   []byte
	StartTime time.Time
}

func (b *BaseContext) AfterAction() {
	body := b.ResponseBody()
	b.DebugLog([]byte("Response"), body)

	resource.TraceLogger.Info(
		misc.FormatTraceLog(&misc.TraceLogArgs{
			TraceId:   b.TraceId,
			Point:     []byte("AllTime"),
			StartTime: b.StartTime,
			EndTime:   time.Now(),
			Msg:       body,
		}),
	)
}

func (b *BaseContext) DebugLog(point, msg []byte) {
	resource.AccessLogger.Debug(misc.FormatAccessLog(b.TraceId, point, msg))
}

func (b *BaseContext) InfoLog(point, msg []byte) {
	resource.AccessLogger.Info(misc.FormatAccessLog(b.TraceId, point, msg))
}

func (b *BaseContext) NoticeLog(point, msg []byte) {
	resource.AccessLogger.Notice(misc.FormatAccessLog(b.TraceId, point, msg))
}

func (b *BaseContext) WarningLog(point, msg []byte) {
	resource.AccessLogger.Warning(misc.FormatAccessLog(b.TraceId, point, msg))
}

func (b *BaseContext) ErrorLog(point, msg []byte) {
	resource.AccessLogger.Error(misc.FormatAccessLog(b.TraceId, point, msg))
}

func (b *BaseContext) CriticalLog(point, msg []byte) {
	resource.AccessLogger.Critical(misc.FormatAccessLog(b.TraceId, point, msg))
}

func (b *BaseContext) AlertLog(point, msg []byte) {
	resource.AccessLogger.Alert(misc.FormatAccessLog(b.TraceId, point, msg))
}

func (b *BaseContext) EmergencyLog(point, msg []byte) {
	resource.AccessLogger.Emergency(misc.FormatAccessLog(b.TraceId, point, msg))
}

type BaseController struct {
}

func (b *BaseController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) gcontroller.ActionContext {
	context := &BaseContext{
		BaseContext: gcontroller.NewBaseContext(req, respWriter),

		StartTime: time.Now(),
	}

	if req.ContentLength < MaxAutoParseBodyLen {
		context.ReqRawBody, _ = ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(context.ReqRawBody))
	}

	_ = req.ParseForm()
	context.QueryValues = req.Form
	context.RemoteRealAddr.Ip, context.RemoteRealAddr.Port = b.parseRemoteAddr(req)

	context.TraceId = b.parseTraceId(context)
	context.NoticeLog(
		[]byte("Request from "+context.RemoteRealAddr.Ip+":"+context.RemoteRealAddr.Port),
		[]byte(req.RequestURI))

	context.ResponseWriter().Header().Add("X-Powered-By", "gohttp")

	return context
}

func (b *BaseController) parseRemoteAddr(req *http.Request) (string, string) {
	rs := strings.Split(req.RemoteAddr, ":")
	if rs[0] == DownstreamServerIp {
		ip := strings.TrimSpace(req.Header.Get(RemoteRealIpHeaderKey))
		port := strings.TrimSpace(req.Header.Get(RemoteRealPortHeaderKey))
		if ip != "" && port != "" {
			return ip, port
		}
	}

	return rs[0], rs[1]
}

func (b *BaseController) parseTraceId(context *BaseContext) []byte {
	traceId := strings.TrimSpace(context.Request().Header.Get(TraceIdHeaderKey))
	if len(traceId) != 0 {
		return []byte(traceId)
	}

	traceId = strings.TrimSpace(context.QueryValues.Get(TraceIdQueryKey))
	if len(traceId) != 0 {
		return []byte(traceId)
	}

	traceIdBytes, _ := misc.DefaultTraceIDGenter.GenID(context.RemoteRealAddr.Ip, context.RemoteRealAddr.Port)
	return traceIdBytes
}
