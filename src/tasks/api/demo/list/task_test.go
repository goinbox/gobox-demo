package list

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/goinbox/mysql"
	"github.com/goinbox/taskflow/v2"

	"gdemo/misc"
	"gdemo/model/demo"
	"gdemo/pcontext"
	"gdemo/test"
)

var (
	ctx      *pcontext.Context
	flowTask *Task
	runner   *taskflow.Runner[*pcontext.Context]

	taskIn  *TaskIn
	taskOut = new(TaskOut)
)

func init() {
	dir, _ := os.Getwd()
	for i := 0; i < 5; i++ {
		dir = filepath.Dir(dir)
	}

	test.InitTestResource(dir)

	ctx = test.Context()
	status := demo.StatusOnline
	taskIn = &TaskIn{
		IDs:        []int64{1, 2, 21},
		Status:     &status,
		ListParams: misc.NewDefaultCommonListParams(),
		ExtSqlQueryItems: []*mysql.SqlColQueryItem{{
			Name:      demo.ColumnName,
			Condition: mysql.SqlCondEqual,
			Value:     "demo",
		}},
	}

	flowTask = NewTask()
	_ = flowTask.Init(taskIn, taskOut)

	runner = taskflow.NewRunner[*pcontext.Context]()
}

func TestMain(m *testing.M) {
	fmt.Println("=== setup")
	setup()

	code := m.Run()

	fmt.Println("=== teardown")
	teardown()

	fmt.Println("=== code", code)
	os.Exit(code)
}

func setup() {
}

func teardown() {
}

func TestRun(t *testing.T) {
	_ = runner.RunTask(ctx, flowTask, taskIn, taskOut)

	t.Log(flowTask.Error(), taskOut)
	t.Log(runner.TaskGraphRunSteps(flowTask, runner.TaskRunSteps()))
}
