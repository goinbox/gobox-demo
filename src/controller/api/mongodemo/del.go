package mongodemo

import (
	"gdemo/errno"

	"github.com/goinbox/exception"
	"github.com/goinbox/gohttp/query"
)

type delActionParams struct {
	id int64
}

func (d *MongoDemoController) DelAction(context *MongoDemoContext) {
	ap, e := d.parseDelActionParams(context)
	if e != nil {
		context.ApiData.Err = e
		return
	}

	deleted, err := context.demoSvc.DeleteById(ap.id)
	if err != nil {
		context.ApiData.Err = exception.New(errno.E_SYS_MONGO_ERROR, err.Error())
		return
	}

	context.ApiData.Data = deleted
}

func (d *MongoDemoController) parseDelActionParams(context *MongoDemoContext) (*delActionParams, *exception.Exception) {
	ap := new(delActionParams)

	qs := query.NewQuerySet()
	qs.Int64Var(&ap.id, "id", true, errno.E_COMMON_INVALID_ARG, "invalid id", query.CheckInt64IsPositive)
	e := qs.Parse(context.QueryValues)

	return ap, e
}
