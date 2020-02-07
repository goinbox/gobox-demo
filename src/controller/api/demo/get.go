package demo

import (
	"gdemo/errno"

	"github.com/goinbox/exception"
	"github.com/goinbox/gohttp/query"
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
		context.ApiData.Err = exception.New(errno.ESysMysqlError, err.Error())
		return
	}

	context.ApiData.Data = entity
}

func (d *DemoController) parseGetActionParams(context *DemoContext) (*getActionParams, *exception.Exception) {
	ap := new(getActionParams)

	qs := query.NewQuerySet()
	qs.Int64Var(&ap.id, "id", true, errno.ECommonInvalidArg, "invalid id", query.CheckInt64IsPositive)
	e := qs.Parse(context.QueryValues)

	return ap, e
}
