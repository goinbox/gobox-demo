package controller

import (
	"gdemo/resource"

	"github.com/goinbox/encoding"
	"github.com/goinbox/gohttp/controller"
	"github.com/goinbox/golog"
	"github.com/goinbox/gomisc"

	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
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
	msg := gomisc.AppendBytes([]byte("[response]\t"), b.RespBody)
	b.AccessLogger.Debug(msg)
}

func (b *BaseContext) Destruct() {
	b.RespBody = nil
	b.QueryValues = nil

	b.AccessLogger.Free()
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

	raddr := []byte(context.RemoteRealAddr.Ip + ":" + context.RemoteRealAddr.Port)
	now := time.Now()
	timeInt := now.UnixNano()
	randInt := gomisc.RandByTime(&now)

	context.TraceId = encoding.Base64Encode(gomisc.AppendBytes(
		raddr, []byte(","),
		[]byte(strconv.FormatInt(timeInt, 10)+","),
		[]byte(strconv.FormatInt(randInt, 10))))

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
