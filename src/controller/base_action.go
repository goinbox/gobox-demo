package controller

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/goinbox/gohttp/httpserver"
	"github.com/goinbox/golog"
	"github.com/google/uuid"

	"gdemo/pcontext"
	"gdemo/resource"
)

const (
	RemoteRealIpHeaderKey   = "REMOTE-REAL-IP"
	RemoteRealPortHeaderKey = "REMOTE-REAL-PORT"

	DownstreamServerIp = "127.0.0.1"

	MaxAutoParseBodyLen = 4096 * 1024

	TraceIDHeaderKey = "TRACE-ID"
	TraceIDQueryKey  = "tid"
)

type BaseAction struct {
	*httpserver.BaseAction

	ReqRawBody []byte

	RemoteRealAddr struct {
		IP   string
		Port int
	}

	Ctx       *pcontext.Context
	StartTime time.Time
}

func NewBaseAction(r *http.Request, w http.ResponseWriter, args []string) *BaseAction {
	a := &BaseAction{
		BaseAction: httpserver.NewBaseAction(r, w, args),
	}

	if r.ContentLength < MaxAutoParseBodyLen {
		a.ReqRawBody, _ = ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(a.ReqRawBody))
	}

	a.RemoteRealAddr.IP, a.RemoteRealAddr.Port = a.parseRemoteAddr(r)

	tid := a.parseTraceID()
	logger := resource.AccessLogger.With(&golog.Field{
		Key:   "tid",
		Value: tid,
	})

	logger.Info("Receive request", []*golog.Field{
		{
			Key:   "ip",
			Value: a.RemoteRealAddr.IP,
		},
		{
			Key:   "port",
			Value: a.RemoteRealAddr.Port,
		},
		{
			Key:   "uri",
			Value: r.RequestURI,
		},
		{
			Key:   "body",
			Value: string(a.ReqRawBody),
		},
	}...)

	a.Ctx = &pcontext.Context{
		TraceID: tid,
		Logger:  logger,
	}

	return a
}

func (a *BaseAction) parseRemoteAddr(r *http.Request) (string, int) {
	rs := strings.Split(r.RemoteAddr, ":")
	if rs[0] == DownstreamServerIp {
		ip := strings.TrimSpace(r.Header.Get(RemoteRealIpHeaderKey))
		portStr := strings.TrimSpace(r.Header.Get(RemoteRealPortHeaderKey))
		if ip != "" && portStr != "" {
			port, _ := strconv.Atoi(portStr)
			return ip, port
		}
	}

	port, _ := strconv.Atoi(rs[1])
	return rs[0], port
}

func (a *BaseAction) parseTraceID() string {
	r := a.Request()

	traceID := strings.TrimSpace(r.Header.Get(TraceIDHeaderKey))
	if len(traceID) != 0 {
		return traceID
	}

	traceID = strings.TrimSpace(r.URL.Query().Get(TraceIDQueryKey))
	if len(traceID) != 0 {
		return traceID
	}

	traceID = uuid.New().String()

	return traceID
}

func (a *BaseAction) Destruct() {
	a.Ctx.Logger.Info("AfterAction", []*golog.Field{
		{
			Key:   "response",
			Value: string(a.ResponseBody()),
		},
		{
			Key:   "cost",
			Value: time.Since(a.StartTime),
		},
	}...)

	a.BaseAction.After()
}
