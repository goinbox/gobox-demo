package del

import (
	"gdemo/pcontext"
	"gdemo/tasks"
)

type TaskIn struct {
	IDs []int64
}

type TaskOut struct {
	RowsAffected int64
}

type Task struct {
	*tasks.BaseTask

	in  *TaskIn
	out *TaskOut

	data struct {
	}
}

func NewTask(ctx *pcontext.Context) *Task {
	t := &Task{
		BaseTask: tasks.NewBaseTask(ctx),
	}

	return t
}

func (t *Task) Name() string {
	return "api.demo.del"
}

func (t *Task) Init(in, out interface{}) error {
	t.in = in.(*TaskIn)
	t.out = out.(*TaskOut)

	return nil
}
