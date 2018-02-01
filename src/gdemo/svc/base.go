package svc

import (
	"github.com/goinbox/golog"
)

type BaseSvc struct {
	elogger golog.ILogger
}

func NewBaseSvc(elogger golog.ILogger) *BaseSvc {
	if elogger == nil {
		elogger = new(golog.NoopLogger)
	}

	return &BaseSvc{
		elogger: elogger,
	}
}
