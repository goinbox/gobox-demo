package del

import (
	"fmt"

	"gdemo/pcontext"
	"github.com/goinbox/taskflow/v2"

	"gdemo/model/factory"
)

func (t *Task) deleteFromDB(ctx *pcontext.Context) (string, error) {
	result := factory.DefaultDaoFactory.DemoDao(ctx).DeleteByIDs(t.in.IDs...)
	if result.Err != nil {
		return "", fmt.Errorf("DemoDao.DeleteByIDs error: %w", result.Err)
	}

	t.out.RowsAffected = result.RowsAffected

	return taskflow.StepCodeSuccess, nil
}
