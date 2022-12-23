package del

import (
	"os"
	"path/filepath"
	"testing"

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
		IDs: []int64{1, 2, 3},
	}, taskOut)
}

func TestDeleteFromDB(t *testing.T) {
	code, err := taskItem.deleteFromDB()
	t.Log(code, err, taskItem.out.RowsAffected)
}
