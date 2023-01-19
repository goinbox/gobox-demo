package tasks

import (
	"fmt"

	"github.com/goinbox/golog"
	"github.com/goinbox/taskflow"

	"gdemo/pcontext"
)

type BaseTask struct {
	err error

	ctx    *pcontext.Context
	logger golog.Logger
}

func NewBaseTask(ctx *pcontext.Context) *BaseTask {
	var logger golog.Logger
	if ctx != nil {
		logger = ctx.Logger
	}

	return &BaseTask{
		ctx:    ctx,
		logger: logger,
	}
}

func (t *BaseTask) Context() *pcontext.Context {
	return t.ctx
}

func (t *BaseTask) Logger() golog.Logger {
	return t.logger
}

func (t *BaseTask) BeforeStep(stepKey string) {
	t.logger = t.ctx.Logger.With(&golog.Field{
		Key:   taskflow.LogFieldKeyStepKey,
		Value: stepKey,
	})
}

func (t *BaseTask) AfterStep(stepKey string) {
	t.logger = t.ctx.Logger
}

func (t *BaseTask) Error() error {
	return t.err
}

func (t *BaseTask) StepFailed(stepKey string, err error) {
	t.err = fmt.Errorf("step %s run error: %w", stepKey, err)
}
