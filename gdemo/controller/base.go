package controller

import (
	"gdemo/idgen"
	"gdemo/misc"
	"gdemo/resource"

	"github.com/goinbox/gohttp/controller"

	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	REMOTE_REAL_IP_HEADER_KEY   = "REMOTE-REAL-IP"
	REMOTE_REAL_PORT_HEADER_KEY = "REMOTE-REAL-PORT"

	DOWNSTREAM_SERVER_IP = "127.0.0.1"

	MAX_AUTO_PARSE_BODY_LEN = 1024 * 1024

	TRACE_ID_HEADER_KEY = "TRACE-ID"
	TRACE_ID_QUERY_KEY  = "traceId"
)

type BaseContext struct {
	Req        *http.Request
	ReqRawBody []byte
	RespWriter http.ResponseWriter
	RespBody   []byte

	QueryValues    url.Values
	RemoteRealAddr struct {
		Ip   string
		Port string
	}

	TraceId   []byte
	StartTime time.Time
}

func (b *BaseContext) Request() *http.Request {
	return b.Req
}

func (b *BaseContext) ResponseWriter() http.ResponseWriter {
	return b.RespWriter
}

func (b *BaseContext) ResponseBody() []byte {
	return b.RespBody
}

func (b *BaseContext) SetResponseBody(body []byte) {
	b.RespBody = body
}

func (b *BaseContext) BeforeAction() {
}

func (b *BaseContext) AfterAction() {
	b.DebugLog([]byte("Response"), b.RespBody)

	resource.TraceLogger.Info(
		misc.FormatTraceLog(&misc.TraceLogArgs{
			TraceId:   b.TraceId,
			Point:     []byte("AllTime"),
			StartTime: b.StartTime,
			EndTime:   time.Now(),
			Msg:       b.RespBody,
		}),
	)
}

func (b *BaseContext) Destruct() {
	b.RespBody = nil
	b.QueryValues = nil
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

func (b *BaseController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) controller.ActionContext {
	context := &BaseContext{
		Req:        req,
		RespWriter: respWriter,
		StartTime:  time.Now(),
	}

	if req.ContentLength < MAX_AUTO_PARSE_BODY_LEN {
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

	context.RespWriter.Header().Add("X-Powered-By", "gohttp")

	return context
}

func (b *BaseController) parseRemoteAddr(req *http.Request) (string, string) {
	rs := strings.Split(req.RemoteAddr, ":")
	if rs[0] == DOWNSTREAM_SERVER_IP {
		ip := strings.TrimSpace(req.Header.Get(REMOTE_REAL_IP_HEADER_KEY))
		port := strings.TrimSpace(req.Header.Get(REMOTE_REAL_PORT_HEADER_KEY))
		if ip != "" && port != "" {
			return ip, port
		}
	}

	return rs[0], rs[1]
}

func (b *BaseController) parseTraceId(context *BaseContext) []byte {
	traceId := strings.TrimSpace(context.Req.Header.Get(TRACE_ID_HEADER_KEY))
	if len(traceId) != 0 {
		return []byte(traceId)
	}

	traceId = strings.TrimSpace(context.QueryValues.Get(TRACE_ID_QUERY_KEY))
	if len(traceId) != 0 {
		return []byte(traceId)
	}

	traceIdBytes, _ := idgen.DefaultTraceIdGenter.GenId(context.RemoteRealAddr.Ip, context.RemoteRealAddr.Port)
	return traceIdBytes
}
