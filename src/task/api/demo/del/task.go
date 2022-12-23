package del

import (
	"gdemo/pcontext"
	"gdemo/task"
)

type TaskIn struct {
	IDs []int64
}

type TaskOut struct {
	RowsAffected int64
}

type Task struct {
	*task.BaseTask

	in  *TaskIn
	out *TaskOut

	data struct {
	}
}

func NewTask(ctx *pcontext.Context) *Task {
	t := &Task{
		BaseTask: task.NewBaseTask(ctx),
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
