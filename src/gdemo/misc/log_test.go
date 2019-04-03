package misc

import (
	"gdemo/conf"
	"time"

	"os"
	"testing"
)

func TestFormatLog(t *testing.T) {
	_ = conf.Init(os.Getenv("GOPATH"))

	msg := FormatAccessLog([]byte("test.format.access"), []byte("test format access log"))
	t.Log(string(msg))

	msg = FormatTraceLog(&TraceLogArgs{
		TraceId:   []byte("abc"),
		Point:     []byte("test.format.trace"),
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Millisecond * 10),
		Username:  []byte("gdemo"),
		Appkey:    []byte("gobox"),
		Msg:       []byte("test format trace log"),
	})
	t.Log(string(msg))
}
