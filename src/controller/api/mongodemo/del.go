package mongodemo

import (
	"gdemo/errno"

	"github.com/goinbox/goerror"
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
		context.ApiData.Err = goerror.New(errno.ESysMongoError, err.Error())
		return
	}

	context.ApiData.Data = deleted
}

func (d *MongoDemoController) parseDelActionParams(context *MongoDemoContext) (*delActionParams, *goerror.Error) {
	ap := new(delActionParams)

	qs := query.NewQuerySet()
	qs.Int64Var(&ap.id, "id", true, errno.ECommonInvalidArg, "invalid id", query.CheckInt64IsPositive)
	e := qs.Parse(context.QueryValues)

	return ap, e
}
