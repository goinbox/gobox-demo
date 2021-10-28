package demo

import (
	"gdemo/perror"
	"github.com/goinbox/gohttp/query"

	"gdemo/perror"
)

type getActionParams struct {
	id int64
}

func (d *DemoController) GetAction(context *DemoContext) {
	ap, e := d.parseGetActionParams(context)
	if e != nil {
		context.ApiData.Err = e
		return
	}

	entity, err := context.demoSvc.GetById(ap.id)
	if err != nil {
		context.ApiData.Err = perror.Error(perror.ESysMysqlError, err.Error())
		return
	}

	context.ApiData.Data = entity
}

func (d *DemoController) parseGetActionParams(context *DemoContext) (*getActionParams, *perror.Error) {
	ap := new(getActionParams)

	qs := query.NewQuerySet()
	qs.Int64Var(&ap.id, "id", true, perror.ECommonInvalidArg, "invalid id", query.CheckInt64IsPositive)
	e := qs.Parse(context.QueryValues)

	return ap, e
}
