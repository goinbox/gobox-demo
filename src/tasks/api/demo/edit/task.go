package edit

import (
	"github.com/goinbox/mysql"

	"gdemo/tasks"
)

type TaskIn struct {
	ID           int64
	UpdateParams interface{}
}

type TaskOut struct {
	RowsAffected int64
}

type Task struct {
	*tasks.BaseTask

	in  *TaskIn
	out *TaskOut

	data struct {
		updateColumns []*mysql.SqlUpdateColumn
	}
}

func NewTask() *Task {
	t := &Task{
		BaseTask: tasks.NewBaseTask(),
	}

	return t
}

func (t *Task) Name() string {
	return "api.demo.edit"
}

func (t *Task) Init(in, out interface{}) error {
	t.in = in.(*TaskIn)
	t.out = out.(*TaskOut)

	return nil
}
