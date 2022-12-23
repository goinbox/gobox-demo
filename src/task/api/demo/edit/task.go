package edit

import (
	"github.com/goinbox/mysql"

	"gdemo/pcontext"
	"gdemo/task"
)

type TaskIn struct {
	ID           int64
	UpdateParams interface{}
}

type TaskOut struct {
	RowsAffected int64
}

type Task struct {
	*task.BaseTask

	in  *TaskIn
	out *TaskOut

	data struct {
		updateColumns []*mysql.SqlUpdateColumn
	}
}

func NewTask(ctx *pcontext.Context) *Task {
	t := &Task{
		BaseTask: task.NewBaseTask(ctx),
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
