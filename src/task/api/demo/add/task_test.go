package add

import (
	"testing"

	"github.com/goinbox/taskflow"

	"gdemo/model/demo"
	"gdemo/test"
)

func TestRun(t *testing.T) {
	task := NewTask(test.Context())
	out := new(TaskOut)
	_ = taskflow.NewRunner(test.Logger()).
		RunTask(task, &TaskIn{
			Name:   "demo",
			Status: demo.StatusOnline,
		}, out)

	t.Log(task.Error(), out)
}
