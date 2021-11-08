package pcontext

import "github.com/goinbox/golog"

type Context struct {
	TraceID string
	Logger  golog.Logger
}
