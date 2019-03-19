package controller

import (
	"gdemo/gvalue"

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
	Rid []byte

	LogFormater  golog.IFormater
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
	b.Rid = nil

	b.AccessLogger.Free()
	b.LogFormater = nil
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

	now := time.Now()
	timeInt := now.UnixNano()
	randInt := gomisc.RandByTime(&now)

	ridStr := context.RemoteRealAddr.Ip + ":" + context.RemoteRealAddr.Port + "," + strconv.FormatInt(timeInt, 10) + "," + strconv.FormatInt(randInt, 10)
	context.Rid = encoding.Base64Encode([]byte(ridStr))

	context.LogFormater = golog.NewWebFormater(context.Rid, []byte(context.RemoteRealAddr.Ip))
	context.RespWriter.Header().Add("X-Powered-By", "gohttp")
	context.AccessLogger = gvalue.NewAsyncLogger(gvalue.AccessLogWriter, context.LogFormater)

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
