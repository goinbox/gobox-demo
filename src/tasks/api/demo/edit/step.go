package edit

import (
	"fmt"

	"gdemo/pcontext"
	"github.com/goinbox/taskflow/v2"

	"gdemo/misc"
	"gdemo/model/factory"
)

func (t *Task) makeUpdateColumns(ctx *pcontext.Context) (string, error) {
	t.data.updateColumns = misc.MakeMysqlUpdateColumns(t.in.UpdateParams)
	if len(t.data.updateColumns) == 0 {
		return taskflow.StepCodeJump1, nil
	}

	return taskflow.StepCodeSuccess, nil
}

func (t *Task) updateDB(ctx *pcontext.Context) (string, error) {
	result := factory.DefaultDaoFactory.DemoDao(ctx).UpdateByIDs(t.data.updateColumns, t.in.ID)
	if result.Err != nil {
		return "", fmt.Errorf("DemoDao.UpdateByIDs error: %w", result.Err)
	}

	t.out.RowsAffected = result.RowsAffected

	return taskflow.StepCodeSuccess, nil
}
