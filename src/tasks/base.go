package tasks

import (
	"fmt"

	"github.com/goinbox/golog"
	"github.com/goinbox/taskflow/v2"

	"gdemo/pcontext"
	"gdemo/tracing"
)

type BaseTask struct {
	err error
}

func NewBaseTask() *BaseTask {
	return &BaseTask{}
}

func (t *BaseTask) BeforeStep(stepKey string) {
}

func (t *BaseTask) AfterStep(stepKey string) {
}

func (t *BaseTask) StepFailed(stepKey string, err error) {
	t.err = fmt.Errorf("step %s run error: %w", stepKey, err)
}

func (t *BaseTask) Error() error {
	return t.err
}

func (t *BaseTask) SetError(err error) {
	t.err = err
}

func RunTask(ctx *pcontext.Context, task taskflow.Task[*pcontext.Context], in, out interface{}) (err error) {
	tctx, span := tracing.StartTrace(ctx, fmt.Sprintf("RunTask %s", task.Name()))
	defer func() { span.EndWithError(err) }()

	defer func() {
		if err != nil {
			tctx.Logger().Error(fmt.Sprintf("RunTask %s error", task.Name()), golog.ErrorField(err))
		}
	}()

	err = taskflow.NewRunner[*pcontext.Context]().
		SetStartTraceFunc(tracing.StartTraceForFramework).
		RunTask(tctx, task, in, out)

	if err != nil {
		return err
	}

	return task.Error()
}
