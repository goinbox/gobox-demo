package demo

import (
	"github.com/goinbox/golog"

	"gdemo/controller/query"
	"gdemo/perror"
)

type editActionParams struct {
	id     int64
	name   string
	status int
}

func (c *DemoController) EditAction(ctx *DemoContext) {
	ap, exists, e := c.parseEditActionParams(ctx)
	if e != nil {
		ctx.ApiData.Err = e
		return
	}

	if len(exists) == 1 { //only has id
		return
	}

	updateFields := make(map[string]interface{})
	_, ok := exists["name"]
	if ok {
		updateFields["name"] = ap.name
	}
	_, ok = exists["status"]
	if ok {
		updateFields["status"] = ap.status
	}

	err := c.demoLogic().UpdateByIDs(ctx.Ctx, updateFields, ap.id)
	if err != nil {
		ctx.Ctx.Logger.Error("demoLogic.UpdateByIDs error", &golog.Field{
			Key:   "err",
			Value: err,
		})

		ctx.ApiData.Err = perror.New(perror.ECommonSysError, "edit error")
		return
	}
}

func (c *DemoController) parseEditActionParams(context *DemoContext) (*editActionParams, map[string]bool, *perror.Error) {
	ap := new(editActionParams)

	qs := query.NewQuerySet()
	qs.Int64Var(&ap.id, "id", true, perror.ECommonInvalidArg, "invalid id", query.CheckInt64IsPositive)
	qs.StringVar(&ap.name, "name", false, perror.ECommonInvalidArg, "invalid name", query.CheckStringNotEmpty)
	qs.IntVar(&ap.status, "status", false, perror.ECommonInvalidArg, "invalid status", nil)
	e := qs.Parse(context.QueryValues)
	if e != nil {
		return ap, nil, e
	}

	if ap.status < 0 {
		return ap, nil, perror.New(perror.ECommonInvalidArg, "invalid status")
	}

	return ap, qs.ExistsInfo(), nil
}
