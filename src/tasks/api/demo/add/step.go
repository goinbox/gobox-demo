package add

import (
	"fmt"

	"github.com/goinbox/taskflow"

	"gdemo/model"
	"gdemo/model/demo"
	"gdemo/model/factory"
)

func (t *Task) genEntity() (string, error) {
	t.data.demoEntity = &demo.Entity{
		BaseEntity: model.BaseEntity{},
		Name:       t.in.Name,
		Status:     t.in.Status,
	}

	return taskflow.StepCodeSuccess, nil
}

func (t *Task) saveEntity() (string, error) {
	result := factory.DefaultDaoFactory.DemoDao(t.Context()).Insert(t.data.demoEntity)
	if result.Err != nil {
		return "", fmt.Errorf("demoDao.Insert error: %w", result.Err)
	}

	t.out.ID = result.LastInsertID

	return taskflow.StepCodeSuccess, nil
}
