package edit

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
		ID: 21,
		UpdateParams: &editParams{
			Name:   "demo",
			Status: demo.StatusOnline,
		},
	}, taskOut)
}

func TestMakeUpdateColumns(t *testing.T) {
	code, err := taskItem.makeUpdateColumns()
	t.Log(code, err, taskItem.data.updateColumns)
}

func TestUpdateDB(t *testing.T) {
	_, _ = taskItem.makeUpdateColumns()
	code, err := taskItem.updateDB()
	t.Log(code, err, taskItem.out.RowsAffected)
}
