package misc

import (
	"github.com/goinbox/gomisc"
	"strconv"

	"time"
)

func FormatAccessLog(point, msg []byte) []byte {
	return gomisc.AppendBytes(
		[]byte("["),
		point,
		[]byte("]"),
		[]byte("\t"),
		msg,
	)
}

type TraceLogArgs struct {
	TraceId   []byte
	Point     []byte
	StartTime time.Time
	EndTime   time.Time
	Username  []byte
	Appkey    []byte
	Msg       []byte
}

var tracePlaceholderBytes = []byte("-")

func FormatTraceLog(tla *TraceLogArgs) []byte {
	var traceId, point, username, appkey, msg []byte

	if tla.TraceId != nil {
		traceId = tla.TraceId
	} else {
		traceId = tracePlaceholderBytes
	}
	if tla.Point != nil {
		point = tla.Point
	} else {
		point = tracePlaceholderBytes
	}
	if tla.Username != nil {
		username = tla.Username
	} else {
		username = tracePlaceholderBytes
	}
	if tla.Appkey != nil {
		appkey = tla.Appkey
	} else {
		appkey = tracePlaceholderBytes
	}
	if tla.Msg != nil {
		msg = tla.Msg
	} else {
		msg = tracePlaceholderBytes
	}

	start := tla.StartTime.UnixNano() / 1000000
	end := tla.EndTime.UnixNano() / 1000000
	elapse := end - start

	return gomisc.AppendBytes(
		traceId,
		[]byte("\t"),
		point,
		[]byte("\t"),
		[]byte(strconv.FormatInt(start, 10)),
		[]byte("\t"),
		[]byte(strconv.FormatInt(end, 10)),
		[]byte("\t"),
		[]byte(strconv.FormatInt(elapse, 10)),
		[]byte("\t"),
		username,
		[]byte("\t"),
		appkey,
		[]byte("\t"),
		msg,
	)
}
