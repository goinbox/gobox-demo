package edit

import (
	"testing"

	"github.com/goinbox/taskflow"

	"gdemo/model/demo"
	"gdemo/test"
)

type editParams struct {
	Name   string
	Status int
}

func TestRun(t *testing.T) {
	task := NewTask(test.Context())
	out := new(TaskOut)
	_ = taskflow.NewRunner(test.Logger()).
		RunTask(task, &TaskIn{
			ID: 21,
			UpdateParams: &editParams{
				Name:   "demo",
				Status: demo.StatusOnline,
			},
		}, out)

	t.Log(task.Error(), out)
}
