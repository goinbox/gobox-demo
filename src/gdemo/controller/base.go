package controller

import (
	"gdemo/idgen"
	"gdemo/misc"
	"gdemo/resource"

	"bytes"
	"github.com/goinbox/gohttp/controller"
	"github.com/goinbox/golog"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	REMOTE_REAL_IP_HEADER_KEY   = "REMOTE-REAL-IP"
	REMOTE_REAL_PORT_HEADER_KEY = "REMOTE-REAL-PORT"

	DOWNSTREAM_SERVER_IP = "127.0.0.1"

	MAX_AUTO_PARSE_BODY_LEN = 1024 * 1024
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

	TraceId      []byte
	AccessLogger golog.ILogger
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
	b.AccessLogger.Debug(misc.FormatAccessLog([]byte("response"), b.RespBody))
}

func (b *BaseContext) Destruct() {
	b.RespBody = nil
	b.QueryValues = nil
}

type BaseController struct {
}

func (b *BaseController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) controller.ActionContext {
	context := &BaseContext{
		Req:        req,
		RespWriter: respWriter,
	}

	if req.ContentLength < MAX_AUTO_PARSE_BODY_LEN {
		context.ReqRawBody, _ = ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(context.ReqRawBody))
	}

	req.ParseForm()
	context.QueryValues = req.Form
	context.RemoteRealAddr.Ip, context.RemoteRealAddr.Port = b.parseRemoteAddr(req)

	context.TraceId, _ = idgen.DefaultTraceIdGenter.GenId(context.RemoteRealAddr.Ip, context.RemoteRealAddr.Port)

	raddr := []byte(context.RemoteRealAddr.Ip + ":" + context.RemoteRealAddr.Port)
	context.AccessLogger = resource.NewLogger(resource.AccessLogWriter,
		golog.NewSimpleFormater().SetAddress(raddr).SetTraceId(context.TraceId))

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
