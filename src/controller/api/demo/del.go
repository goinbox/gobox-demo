package demo

import (
	"github.com/goinbox/golog"

	"gdemo/controller/query"
	"gdemo/perror"
)

type delActionParams struct {
	id int64
}

func (c *DemoController) DelAction(ctx *DemoContext) {
	ap, e := c.parseDelActionParams(ctx)
	if e != nil {
		ctx.ApiData.Err = e
		return
	}

	err := c.demoLogic().DeleteByIDs(ctx.Ctx, ap.id)
	if err != nil {
		ctx.Ctx.Logger.Error("demoLogic.DeleteByIDs error", golog.ErrorField(err))

		ctx.ApiData.Err = perror.New(perror.ECommonSysError, "del error")
		return
	}
}

func (c *DemoController) parseDelActionParams(context *DemoContext) (*delActionParams, *perror.Error) {
	ap := new(delActionParams)

	qs := query.NewQuerySet()
	qs.Int64Var(&ap.id, "id", true, perror.ECommonInvalidArg, "invalid id", query.CheckInt64IsPositive)
	e := qs.Parse(context.QueryValues)

	return ap, e
}
