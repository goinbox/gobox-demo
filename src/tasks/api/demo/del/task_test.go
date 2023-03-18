package del

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/goinbox/taskflow"

	"gdemo/pcontext"
	"gdemo/test"
)

var (
	ctx      *pcontext.Context
	flowTask *Task
	runner   *taskflow.Runner

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
	taskIn = &TaskIn{
		IDs: []int64{1, 2, 3},
	}

	flowTask = NewTask(test.Context())
	_ = flowTask.Init(taskIn, taskOut)

	runner = taskflow.NewRunner(ctx.Logger)
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
	_ = runner.RunTask(flowTask, taskIn, taskOut)

	t.Log(flowTask.Error(), taskOut)
	t.Log(runner.TaskGraphRunSteps(flowTask, runner.TaskRunSteps()))
}
