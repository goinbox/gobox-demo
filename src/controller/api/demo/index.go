package demo

import (
	"github.com/goinbox/golog"

	"gdemo/controller/query"
	"gdemo/model/demo"
	"gdemo/perror"

	"github.com/goinbox/mysql"
)

type indexActionRequestParams struct {
	status int

	offset int64
	cnt    int64
}

type indexActionResponseData struct {
	Total    int64
	DemoList []*demo.Entity
}

func (c *DemoController) IndexAction(ctx *DemoContext) {
	ap, exists, e := c.parseIndexActionParams(ctx)
	if e != nil {
		ctx.ApiData.Err = e
		return
	}

	var conds []*mysql.SqlColQueryItem
	if exists["status"] {
		conds = append(conds, &mysql.SqlColQueryItem{
			Name:      "status",
			Condition: mysql.SqlCondEqual,
			Value:     ap.status,
		})
	}

	params := &mysql.SqlQueryParams{
		CondItems: conds,

		OrderBy: "id desc",
		Offset:  ap.offset,
		Cnt:     ap.cnt,
	}

	logic := c.demoLogic()

	entities, err := logic.SimpleQueryAnd(ctx.Ctx, params)
	if err != nil {
		ctx.Ctx.Logger.Error("demoLogic.SimpleQueryAnd error", golog.ErrorField(err))

		ctx.ApiData.Err = perror.New(perror.ECommonSysError, "query error")
		return
	}

	total, err := logic.SimpleTotalAnd(ctx.Ctx, conds...)
	if err != nil {
		ctx.Ctx.Logger.Error("demoLogic.SimpleTotalAnd error", golog.ErrorField(err))

		ctx.ApiData.Err = perror.New(perror.ECommonSysError, "query error")
		return
	}

	ctx.ApiData.Data = &indexActionResponseData{
		Total:    total,
		DemoList: entities,
	}
}

func (c *DemoController) parseIndexActionParams(context *DemoContext) (*indexActionRequestParams, map[string]bool, *perror.Error) {
	ap := new(indexActionRequestParams)

	qs := query.NewQuerySet()
	qs.IntVar(&ap.status, "status", false, perror.ECommonInvalidArg, "invalid status", nil)
	qs.Int64Var(&ap.offset, "offset", false, perror.ECommonInvalidArg, "invalid offset", nil)
	qs.Int64Var(&ap.cnt, "cnt", false, perror.ECommonInvalidArg, "invalid cnt", nil)
	e := qs.Parse(context.QueryValues)
	if e != nil {
		return ap, nil, e
	}

	if ap.status < 0 {
		return ap, nil, perror.New(perror.ECommonInvalidArg, "invalid status")
	}
	if ap.offset < 0 {
		return ap, nil, perror.New(perror.ECommonInvalidArg, "invalid offset")
	}
	if ap.cnt < 0 {
		return ap, nil, perror.New(perror.ECommonInvalidArg, "invalid cnt")
	}

	return ap, qs.ExistsInfo(), nil
}
