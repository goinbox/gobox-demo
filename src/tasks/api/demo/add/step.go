package add

import (
	"fmt"

	"gdemo/pcontext"
	"github.com/goinbox/taskflow/v2"

	"gdemo/model"
	"gdemo/model/demo"
	"gdemo/model/factory"
)

func (t *Task) genEntity(ctx *pcontext.Context) (string, error) {
	t.data.demoEntity = &demo.Entity{
		BaseEntity: model.BaseEntity{},
		Name:       t.in.Name,
		Status:     t.in.Status,
	}

	return taskflow.StepCodeSuccess, nil
}

func (t *Task) saveEntity(ctx *pcontext.Context) (string, error) {
	result := factory.DefaultDaoFactory.DemoDao(ctx).Insert(ctx, t.data.demoEntity)
	if result.Err != nil {
		return "", fmt.Errorf("demoDao.Insert error: %w", result.Err)
	}

	t.out.ID = result.LastInsertID

	return taskflow.StepCodeSuccess, nil
}
