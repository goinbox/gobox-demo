package demo

import (
	"github.com/goinbox/golog"

	"gdemo/controller/query"
	"gdemo/perror"
)

type getActionParams struct {
	id int64
}

func (c *DemoController) GetAction(ctx *DemoContext) {
	ap, e := c.parseGetActionParams(ctx)
	if e != nil {
		ctx.ApiData.Err = e
		return
	}

	entity, err := c.demoLogic().SelectByID(ctx.Ctx, ap.id)
	if err != nil {
		ctx.Ctx.Logger.Error("demoLogic.SelectByID error", &golog.Field{
			Key:   "err",
			Value: err,
		})

		ctx.ApiData.Err = perror.New(perror.ECommonSysError, "get error")
		return
	}

	ctx.ApiData.Data = entity
}

func (c *DemoController) parseGetActionParams(context *DemoContext) (*getActionParams, *perror.Error) {
	ap := new(getActionParams)

	qs := query.NewQuerySet()
	qs.Int64Var(&ap.id, "id", true, perror.ECommonInvalidArg, "invalid id", query.CheckInt64IsPositive)
	e := qs.Parse(context.QueryValues)

	return ap, e
}
