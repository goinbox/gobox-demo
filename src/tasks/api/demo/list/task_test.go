package list

import (
	"testing"

	"github.com/goinbox/taskflow"

	"gdemo/misc"
	"gdemo/model/demo"
	"gdemo/test"
)

func TestRun(t *testing.T) {
	task := NewTask(test.Context())
	out := new(TaskOut)
	status := demo.StatusOnline
	_ = taskflow.NewRunner(test.Logger()).
		RunTask(task, &TaskIn{
			IDs:              []int64{1, 2, 21},
			Status:           &status,
			ListParams:       misc.NewDefaultCommonListParams(),
			ExtSqlQueryItems: nil,
		}, out)

	t.Log(task.Error(), out)
	for _, entity := range out.DemoList {
		t.Log(entity, *entity.ID, *entity.AddTime, *entity.EditTime)
	}
}
