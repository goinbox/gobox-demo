package controller

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	gcontroller "github.com/goinbox/gohttp/controller"
	"github.com/goinbox/golog"

	"gdemo/misc"
	"gdemo/pcontext"
	"gdemo/resource"
)

const (
	RemoteRealIpHeaderKey   = "REMOTE-REAL-IP"
	RemoteRealPortHeaderKey = "REMOTE-REAL-PORT"

	DownstreamServerIp = "127.0.0.1"

	MaxAutoParseBodyLen = 1024 * 1024

	TraceIDHeaderKey = "TRACE-ID"
	TraceIDQueryKey  = "tid"
)

type BaseContext struct {
	*gcontroller.BaseContext

	ReqRawBody []byte

	QueryValues    url.Values
	RemoteRealAddr struct {
		IP   string
		Port int
	}

	Ctx       *pcontext.Context
	StartTime time.Time
}

func (c *BaseContext) AfterAction() {
	c.Ctx.Logger.Info("AfterAction", []*golog.Field{
		{
			Key:   "response",
			Value: string(c.ResponseBody()),
		},
		{
			Key:   "cost",
			Value: time.Since(c.StartTime),
		},
	}...)
}

type BaseController struct {
}

func (c *BaseController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) gcontroller.ActionContext {
	actx := &BaseContext{
		BaseContext: gcontroller.NewBaseContext(req, respWriter),
		StartTime:   time.Now(),
	}

	if req.ContentLength < MaxAutoParseBodyLen {
		actx.ReqRawBody, _ = ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(actx.ReqRawBody))
	}

	_ = req.ParseForm()
	actx.QueryValues = req.Form
	actx.RemoteRealAddr.IP, actx.RemoteRealAddr.Port = c.parseRemoteAddr(req)

	tid := c.parseTraceID(actx)
	logger := resource.AccessLogger.With(&golog.Field{
		Key:   "tid",
		Value: tid,
	})

	logger.Info("Receive request", []*golog.Field{
		{
			Key:   "ip",
			Value: actx.RemoteRealAddr.IP,
		},
		{
			Key:   "port",
			Value: actx.RemoteRealAddr.Port,
		},
		{
			Key:   "uri",
			Value: req.RequestURI,
		},
	}...)

	actx.Ctx = &pcontext.Context{
		TraceID: tid,
		Logger:  logger,
	}

	return actx
}

func (c *BaseController) parseRemoteAddr(req *http.Request) (string, int) {
	rs := strings.Split(req.RemoteAddr, ":")
	if rs[0] == DownstreamServerIp {
		ip := strings.TrimSpace(req.Header.Get(RemoteRealIpHeaderKey))
		portStr := strings.TrimSpace(req.Header.Get(RemoteRealPortHeaderKey))
		if ip != "" && portStr != "" {
			port, _ := strconv.Atoi(portStr)
			return ip, port
		}
	}

	port, _ := strconv.Atoi(rs[1])
	return rs[0], port
}

func (c *BaseController) parseTraceID(context *BaseContext) string {
	traceID := strings.TrimSpace(context.Request().Header.Get(TraceIDHeaderKey))
	if len(traceID) != 0 {
		return traceID
	}

	traceID = strings.TrimSpace(context.QueryValues.Get(TraceIDQueryKey))
	if len(traceID) != 0 {
		return traceID
	}

	traceID, _ = misc.DefaultTraceIDGenter.GenID(context.RemoteRealAddr.IP, context.RemoteRealAddr.Port)
	return traceID
}
