package list

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goinbox/mysql"

	"gdemo/misc"
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
	status := demo.StatusOnline
	_ = taskItem.Init(&TaskIn{
		IDs:        []int64{1, 2, 21},
		Status:     &status,
		ListParams: misc.NewDefaultCommonListParams(),
		ExtSqlQueryItems: []*mysql.SqlColQueryItem{{
			Name:      demo.ColumnName,
			Condition: mysql.SqlCondEqual,
			Value:     "demo",
		}},
	}, taskOut)
}

func TestMakeSqlQueryParams(t *testing.T) {
	code, err := taskItem.makeSqlQueryParams()
	t.Log(code, err, taskItem.data.queryParams)
	for _, item := range taskItem.data.queryParams.CondItems {
		t.Log(item)
	}
}

func TestQueryFromDB(t *testing.T) {
	_, _ = taskItem.makeSqlQueryParams()
	code, err := taskItem.queryFromDB()
	t.Log(code, err, taskItem.out)
	for _, entity := range taskItem.out.DemoList {
		t.Log(entity, *entity.ID, *entity.AddTime, *entity.EditTime)
	}
}
