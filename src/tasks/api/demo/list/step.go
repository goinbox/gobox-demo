package list

import (
	"fmt"

	"gdemo/pcontext"
	"github.com/goinbox/mysql"
	"github.com/goinbox/taskflow/v2"

	"gdemo/model"
	"gdemo/model/demo"
	"gdemo/model/factory"
)

func (t *Task) makeSqlQueryParams(ctx *pcontext.Context) (string, error) {
	var condItems []*mysql.SqlColQueryItem
	item := model.MakeInt64SliceSqlColQueryItem(demo.ColumnID, t.in.IDs)
	if item != nil {
		condItems = append(condItems, item)
	}
	if t.in.Status != nil {
		condItems = append(condItems, &mysql.SqlColQueryItem{
			Name:      demo.ColumnStatus,
			Condition: mysql.SqlCondEqual,
			Value:     *t.in.Status,
		})
	}
	if len(t.in.ExtSqlQueryItems) > 0 {
		condItems = append(condItems, t.in.ExtSqlQueryItems...)
	}

	t.data.queryParams = &mysql.SqlQueryParams{
		CondItems: condItems,
		OrderBy:   t.in.ListParams.MysqlOrderBy(),
		Offset:    t.in.ListParams.Offset,
		Cnt:       t.in.ListParams.Limit,
	}

	return taskflow.StepCodeSuccess, nil
}

func (t *Task) queryFromDB(ctx *pcontext.Context) (string, error) {
	dao := factory.DefaultDaoFactory.DemoDao(ctx)
	total, err := dao.SimpleTotalAnd(ctx, t.data.queryParams.CondItems...)
	if err != nil {
		return "", fmt.Errorf("dao.SimpleTotalAnd error: %w", err)
	}
	if total > 0 {
		err = dao.SimpleQueryAnd(ctx, t.data.queryParams, &t.out.DemoList)
		if err != nil {
			return "", fmt.Errorf("dao.SimpleQueryAnd error: %w", err)
		}
	}

	t.out.Total = total

	return taskflow.StepCodeSuccess, nil
}
