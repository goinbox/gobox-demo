package add

import (
	"gdemo/model/demo"
	"gdemo/model/factory"
	"gdemo/pcontext"
	"gdemo/task"
)

type TaskIn struct {
	Name   string
	Status int
}

type TaskOut struct {
	ID int64
}

type Task struct {
	*task.BaseTask

	in  *TaskIn
	out *TaskOut

	demoDao demo.Dao

	data struct {
		demoEntity *demo.Entity
	}
}

func NewTask(ctx *pcontext.Context) *Task {
	t := &Task{
		BaseTask: task.NewBaseTask(ctx),
	}

	return t
}

func (t *Task) Name() string {
	return "api.demo.add"
}

func (t *Task) Init(in, out interface{}) error {
	t.in = in.(*TaskIn)
	t.out = out.(*TaskOut)

	ctx := t.Context()
	t.demoDao = factory.DefaultDaoFactory.DemoDao(ctx)

	return nil
}
