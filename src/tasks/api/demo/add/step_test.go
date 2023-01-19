package add

import (
	"os"
	"path/filepath"
	"testing"

	"gdemo/model/demo"
	"gdemo/test"
)

var (
	taskItem *Task
	taskOut  = new(TaskOut)
)

func init() {
	dir, _ := os.Getwd()
	for i := 0; i < 5; i++ {
		dir = filepath.Dir(dir)
	}

	test.InitTestResource(dir)

	taskItem = NewTask(test.Context())
	_ = taskItem.Init(&TaskIn{
		Name:   "demo",
		Status: demo.StatusOnline,
	}, taskOut)
}

func TestGenEntity(t *testing.T) {
	code, err := taskItem.genEntity()
	t.Log(code, err, taskItem.data.demoEntity)
}

func TestSaveEntity(t *testing.T) {
	_, _ = taskItem.genEntity()
	code, err := taskItem.saveEntity()
	t.Log(code, err, taskItem.out)
}
