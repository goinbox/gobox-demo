package demo

import (
	"github.com/goinbox/golog"

	"gdemo/controller/query"
	"gdemo/model"
	"gdemo/model/demo"
	"gdemo/perror"
)

type addActionParams struct {
	name   string
	status int
}

func (c *DemoController) AddAction(actx *DemoContext) {
	ap, e := c.parseAddActionParams(actx)
	if e != nil {
		actx.ApiData.Err = e
		return
	}

	err := c.demoLogic().Insert(actx.Ctx, &demo.Entity{
		BaseEntity: model.BaseEntity{},
		Name:       ap.name,
		Status:     ap.status,
	})
	if err != nil {
		if model.DuplicateError(err) {
			actx.ApiData.Err = perror.New(perror.ECommonDataAlreadyExist, "data already exist")
			return
		}

		actx.Ctx.Logger.Error("demoLogic.Insert error", golog.ErrorField(err))

		actx.ApiData.Err = perror.New(perror.ECommonSysError, "add error")
		return
	}
}

func (c *DemoController) parseAddActionParams(context *DemoContext) (*addActionParams, *perror.Error) {
	ap := new(addActionParams)

	qs := query.NewQuerySet()
	qs.StringVar(&ap.name, "name", true, perror.ECommonInvalidArg, "invalid name", query.CheckStringNotEmpty)
	qs.IntVar(&ap.status, "status", true, perror.ECommonInvalidArg, "invalid status", nil)
	e := qs.Parse(context.QueryValues)
	if e != nil {
		return ap, e
	}

	if ap.status < 0 {
		return ap, perror.New(perror.ECommonInvalidArg, "invalid status")
	}

	return ap, nil
}
