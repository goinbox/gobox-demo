package mongodemo

import (
	"gdemo/errno"

	"github.com/goinbox/exception"
	"github.com/goinbox/gohttp/query"
)

type getActionParams struct {
	id int64
}

func (d *MongoDemoController) GetAction(context *MongoDemoContext) {
	ap, e := d.parseGetActionParams(context)
	if e != nil {
		context.ApiData.Err = e
		return
	}

	entity, err := context.demoSvc.GetById(ap.id)
	if err != nil {
		context.ApiData.Err = exception.New(errno.ESysMongoError, err.Error())
		return
	}
	context.ApiData.Data = entity
}

func (d *MongoDemoController) parseGetActionParams(context *MongoDemoContext) (*getActionParams, *exception.Exception) {
	ap := new(getActionParams)

	qs := query.NewQuerySet()
	qs.Int64Var(&ap.id, "id", true, errno.ECommonInvalidArg, "invalid id", query.CheckInt64IsPositive)
	e := qs.Parse(context.QueryValues)

	return ap, e
}
