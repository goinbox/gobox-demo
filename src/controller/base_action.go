package controller

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/goinbox/gohttp/v6/httpserver"
	"github.com/goinbox/golog"
	"github.com/google/uuid"

	"gdemo/pcontext"
	"gdemo/resource"
)

const (
	remoteRealIpHeaderKey   = "REMOTE-REAL-IP"
	remoteRealPortHeaderKey = "REMOTE-REAL-PORT"

	downstreamServerIp = "127.0.0.1"

	maxAutoParseBodyLen = 4096 * 1024

	traceIDHeaderKey = "TRACE-ID"
	traceIDQueryKey  = "tid"
)

type BaseAction struct {
	httpserver.BaseAction[*pcontext.Context]

	ReqRawBody []byte

	RemoteRealAddr struct {
		IP   string
		Port int
	}

	StartTime time.Time
}

func (a *BaseAction) Init(r *http.Request, w http.ResponseWriter, args []string) *pcontext.Context {
	a.BaseAction.Init(r, w, args)

	if r.ContentLength < maxAutoParseBodyLen {
		a.ReqRawBody, _ = io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(a.ReqRawBody))
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

	return pcontext.NewContext(logger).SetTraceID(tid)
}

func (a *BaseAction) parseRemoteAddr(r *http.Request) (string, int) {
	rs := strings.Split(r.RemoteAddr, ":")
	if rs[0] == downstreamServerIp {
		ip := strings.TrimSpace(r.Header.Get(remoteRealIpHeaderKey))
		portStr := strings.TrimSpace(r.Header.Get(remoteRealPortHeaderKey))
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

	traceID := strings.TrimSpace(r.Header.Get(traceIDHeaderKey))
	if len(traceID) != 0 {
		return traceID
	}

	traceID = strings.TrimSpace(r.URL.Query().Get(traceIDQueryKey))
	if len(traceID) != 0 {
		return traceID
	}

	traceID = uuid.New().String()

	return traceID
}

func (a *BaseAction) Before(ctx *pcontext.Context) error {
	_ = a.BaseAction.Before(ctx)

	a.StartTime = time.Now()

	return nil
}

func (a *BaseAction) After(ctx *pcontext.Context, err error) {
	ctx.Logger().Info("AfterAction", []*golog.Field{
		{
			Key:   "response",
			Value: string(a.ResponseBody()),
		},
		{
			Key:   "cost",
			Value: time.Since(a.StartTime),
		},
	}...)

	a.BaseAction.After(ctx, err)
}
