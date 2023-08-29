package list

import (
	"github.com/goinbox/mysql"

	"gdemo/misc"
	"gdemo/model/demo"
	"gdemo/tasks"
)

type TaskIn struct {
	IDs    []int64
	Status *int

	ListParams       *misc.CommonListParams
	ExtSqlQueryItems []*mysql.SqlColQueryItem
}

type TaskOut struct {
	Total    int64
	DemoList []*demo.Entity
}

type Task struct {
	*tasks.BaseTask

	in  *TaskIn
	out *TaskOut

	data struct {
		queryParams *mysql.SqlQueryParams
	}
}

func NewTask() *Task {
	t := &Task{
		BaseTask: tasks.NewBaseTask(),
	}

	return t
}

func (t *Task) Name() string {
	return "api.demo.list"
}

func (t *Task) Init(in, out interface{}) error {
	t.in = in.(*TaskIn)
	t.out = out.(*TaskOut)

	return nil
}
