package api

import (
	"fmt"

	"github.com/goinbox/golog"
	"github.com/goinbox/taskflow"

	"gdemo/pcontext"
)

func RunTask(ctx *pcontext.Context, task taskflow.Task, in, out interface{}) error {
	err := taskflow.NewRunner(ctx.Logger).RunTask(task, in, out)
	if err != nil {
		ctx.Logger.Error(fmt.Sprintf("Runner.RunTask %s error", task.Name()), golog.ErrorField(err))
		return err
	}

	err = task.Error()
	if err != nil {
		ctx.Logger.Error(fmt.Sprintf("RunTask %s end with error", task.Name()), golog.ErrorField(err))
	}

	return err
}
