package edit

import (
	"fmt"

	"github.com/goinbox/taskflow"

	"gdemo/misc"
	"gdemo/model/factory"
)

func (t *Task) makeUpdateColumns() (string, error) {
	t.data.updateColumns = misc.MakeMysqlUpdateColumns(t.in.UpdateParams)
	if len(t.data.updateColumns) == 0 {
		return taskflow.StepCodeJump1, nil
	}

	return taskflow.StepCodeSuccess, nil
}

func (t *Task) updateDB() (string, error) {
	result := factory.DefaultDaoFactory.DemoDao(t.Context()).UpdateByIDs(t.data.updateColumns, t.in.ID)
	if result.Err != nil {
		return "", fmt.Errorf("DemoDao.UpdateByIDs error: %w", result.Err)
	}

	t.out.RowsAffected = result.RowsAffected

	return taskflow.StepCodeSuccess, nil
}
