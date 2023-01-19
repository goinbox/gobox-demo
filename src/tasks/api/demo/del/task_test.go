package del

import (
	"testing"

	"github.com/goinbox/taskflow"

	"gdemo/test"
)

func TestRun(t *testing.T) {
	task := NewTask(test.Context())
	out := new(TaskOut)
	_ = taskflow.NewRunner(test.Logger()).
		RunTask(task, &TaskIn{
			IDs: []int64{1, 2, 3},
		}, out)

	t.Log(task.Error(), out)
}
