package misc

import (
	"github.com/goinbox/gomisc"

	"strconv"
	"time"
)

func FormatAccessLog(traceid, point, msg []byte) []byte {
	return gomisc.AppendBytes(
		traceid, []byte("\t"),
		[]byte("["), point, []byte("]"), []byte("\t"),
		msg,
	)
}

type TraceLogArgs struct {
	TraceId   []byte
	Point     []byte
	StartTime time.Time
	EndTime   time.Time
	Msg       []byte
}

var tracePlaceholderBytes = []byte("-")

func FormatTraceLog(tla *TraceLogArgs) []byte {
	var traceId, point, msg []byte

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
		msg,
	)
}
