package add

import (
	"gdemo/model/demo"
	"gdemo/pcontext"
	"gdemo/tasks"
)

type TaskIn struct {
	Name   string
	Status int
}

type TaskOut struct {
	ID int64
}

type Task struct {
	*tasks.BaseTask

	in  *TaskIn
	out *TaskOut

	data struct {
		demoEntity *demo.Entity
	}
}

func NewTask(ctx *pcontext.Context) *Task {
	t := &Task{
		BaseTask: tasks.NewBaseTask(ctx),
	}

	return t
}

func (t *Task) Name() string {
	return "api.demo.add"
}

func (t *Task) Init(in, out interface{}) error {
	t.in = in.(*TaskIn)
	t.out = out.(*TaskOut)

	return nil
}
