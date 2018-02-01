package demo

import (
	"gdemo/errno"

	"github.com/goinbox/exception"
	"github.com/goinbox/gohttp/query"
)

type getActionParams struct {
	id int64
}

func (this *DemoController) GetAction(context *DemoContext) {
	ap, e := this.parseGetActionParams(context)
	if e != nil {
		context.ApiData.Err = e
		return
	}

	entity, err := context.demoSvc.GetById(ap.id)
	if err != nil {
		context.ApiData.Err = exception.New(errno.E_SYS_MYSQL_ERROR, err.Error())
		return
	}

	context.ApiData.Data = entity
}

func (this *DemoController) parseGetActionParams(context *DemoContext) (*getActionParams, *exception.Exception) {
	ap := new(getActionParams)

	qs := query.NewQuerySet()
	qs.Int64Var(&ap.id, "id", true, errno.E_COMMON_INVALID_ID, "invalid id", query.CheckInt64IsPositive)
	e := qs.Parse(context.QueryValues)

	return ap, e
}
