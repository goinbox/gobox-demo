package del

import (
	"fmt"

	"github.com/goinbox/taskflow"

	"gdemo/model/factory"
)

func (t *Task) deleteFromDB() (string, error) {
	result := factory.DefaultDaoFactory.DemoDao(t.Context()).DeleteByIDs(t.in.IDs...)
	if result.Err != nil {
		return "", fmt.Errorf("DemoDao.DeleteByIDs error: %w", result.Err)
	}

	t.out.RowsAffected = result.RowsAffected

	return taskflow.StepCodeSuccess, nil
}
