package mongo

import (
	"github.com/goinbox/golog"
)

type mongoLogger struct {
	logger golog.ILogger
}

func NewMongoLogger(logger golog.ILogger) *mongoLogger {
	this := &mongoLogger{
		logger: logger,
	}
	return this
}

func (this *mongoLogger) Output(calldepth int, s string) (err error) {
	this.logger.Log(calldepth, []byte(s))
	return
}
