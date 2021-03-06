package demo

import (
	"github.com/goinbox/goerror"
	"github.com/goinbox/gohttp/query"

	"gdemo/errno"
)

type delActionParams struct {
	id int64
}

func (d *DemoController) DelAction(context *DemoContext) {
	ap, e := d.parseDelActionParams(context)
	if e != nil {
		context.ApiData.Err = e
		return
	}

	deleted, err := context.demoSvc.DeleteById(ap.id)
	if err != nil {
		context.ApiData.Err = goerror.New(errno.ESysMysqlError, err.Error())
		return
	}

	context.ApiData.Data = deleted
}

func (d *DemoController) parseDelActionParams(context *DemoContext) (*delActionParams, *goerror.Error) {
	ap := new(delActionParams)

	qs := query.NewQuerySet()
	qs.Int64Var(&ap.id, "id", true, errno.ECommonInvalidArg, "invalid id", query.CheckInt64IsPositive)
	e := qs.Parse(context.QueryValues)

	return ap, e
}
