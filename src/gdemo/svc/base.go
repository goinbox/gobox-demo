package svc

import (
	"github.com/goinbox/golog"
)

type BaseSvc struct {
	Elogger golog.ILogger
}

func NewBaseSvc(elogger golog.ILogger) *BaseSvc {
	if elogger == nil {
		elogger = new(golog.NoopLogger)
	}

	return &BaseSvc{
		Elogger: elogger,
	}
}
