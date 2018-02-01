package demo

import (
	"gdemo/errno"

	"github.com/goinbox/exception"
	"github.com/goinbox/gohttp/query"
)

type delActionParams struct {
	id int64
}

func (this *DemoController) DelAction(context *DemoContext) {
	ap, e := this.parseDelActionParams(context)
	if e != nil {
		context.ApiData.Err = e
		return
	}

	deleted, err := context.demoSvc.DeleteById(ap.id)
	if err != nil {
		context.ApiData.Err = exception.New(errno.E_SYS_MYSQL_ERROR, err.Error())
		return
	}

	context.ApiData.Data = deleted
}

func (this *DemoController) parseDelActionParams(context *DemoContext) (*delActionParams, *exception.Exception) {
	ap := new(delActionParams)

	qs := query.NewQuerySet()
	qs.Int64Var(&ap.id, "id", true, errno.E_COMMON_INVALID_ID, "invalid id", query.CheckInt64IsPositive)
	e := qs.Parse(context.QueryValues)

	return ap, e
}
