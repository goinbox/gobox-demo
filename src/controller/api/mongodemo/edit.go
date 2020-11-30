package mongodemo

import (
	"gdemo/errno"
	demoSvc "gdemo/svc/mongodemo"

	"github.com/goinbox/goerror"
	"github.com/goinbox/gohttp/query"
)

type editActionParams struct {
	entiry demoSvc.MongoDemoEntity

	tid int64
}

func (d *MongoDemoController) EditAction(context *MongoDemoContext) {
	ap, exists, e := d.parseEditActionParams(context)
	if e != nil {
		context.ApiData.Err = e
		return
	}

	if len(exists) == 1 { //only has id
		return
	}

	updated, err := context.demoSvc.UpdateById(ap.Id, ap, exists)
	if err != nil {
		context.ApiData.Err = goerror.New(errno.ECommonUpdateEntityFailed, err.Error())
		return
	}

	context.ApiData.Data = updated
}

func (d *MongoDemoController) parseEditActionParams(context *MongoDemoContext) (*demoSvc.MongoDemoEntity, map[string]bool, *goerror.Error) {
	ap := new(editActionParams)

	qs := query.NewQuerySet()

	qs.Int64Var(&ap.tid, "id", true, errno.ECommonInvalidArg, "invalid id", query.CheckInt64IsPositive)
	qs.StringVar(&ap.entiry.Name, "name", false, errno.ECommonInvalidArg, "invalid name", query.CheckStringNotEmpty)
	qs.IntVar(&ap.entiry.Status, "status", false, errno.ECommonInvalidArg, "invalid status", nil)
	e := qs.Parse(context.QueryValues)

	if e != nil {
		return &ap.entiry, nil, e
	}

	ap.entiry.Id = ap.tid

	if ap.entiry.Status < 0 {
		return &ap.entiry, nil, goerror.New(errno.ECommonInvalidArg, "invalid status")
	}

	return &ap.entiry, qs.ExistsInfo(), nil
}
