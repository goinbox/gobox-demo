package pcontext

import (
	"github.com/goinbox/golog"
	"github.com/goinbox/mysql"
)

type Context struct {
	TraceID string
	Logger  golog.Logger

	MySQLClient *mysql.Client
}
