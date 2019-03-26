/**
* @file base.go
* @brief format msg before send to writer
* @author ligang
* @date 2016-07-12
 */

package golog

import (
	"github.com/goinbox/gomisc"

	"time"
)

type simpleFormater struct {
	address []byte
	traceId []byte

	timeLayout string
}

func NewSimpleFormater() *simpleFormater {
	return &simpleFormater{
		[]byte("-"),
		[]byte("-"),
		gomisc.TimeGeneralLayout(),
	}
}

func (s *simpleFormater) SetAddress(address []byte) *simpleFormater {
	s.address = address

	return s
}

func (s *simpleFormater) SetTraceId(traceId []byte) *simpleFormater {
	s.traceId = traceId

	return s
}

func (s *simpleFormater) SetTimeLayout(layout string) *simpleFormater {
	s.timeLayout = layout

	return s
}

func (s *simpleFormater) Format(level int, msg []byte) []byte {
	lm, ok := LogLevels[level]
	if !ok {
		lm = []byte("-")
	}

	return gomisc.AppendBytes(
		[]byte("["), lm, []byte("]\t"),
		[]byte("["), []byte(time.Now().Format(s.timeLayout)), []byte("]\t"),
		s.address, []byte("\t"),
		s.traceId, []byte("\t"),
		msg,
	)
}
