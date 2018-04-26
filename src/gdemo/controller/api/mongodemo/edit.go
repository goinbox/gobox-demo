package mongodemo

import (
	"gdemo/errno"
	demoSvc "gdemo/svc/mongodemo"

	"github.com/goinbox/exception"
	"github.com/goinbox/gohttp/query"
)

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
		context.ApiData.Err = exception.New(errno.E_API_DEMO_UPDATE_FAILED, err.Error())
		return
	}

	context.ApiData.Data = updated
}

func (d *MongoDemoController) parseEditActionParams(context *MongoDemoContext) (*demoSvc.MongoDemoEntity, map[string]bool, *exception.Exception) {
	ap := new(demoSvc.MongoDemoEntity)

	qs := query.NewQuerySet()
	qs.Int64Var(&ap.Id, "id", true, errno.E_COMMON_INVALID_ID, "invalid id", query.CheckInt64IsPositive)
	qs.StringVar(&ap.Name, "name", false, errno.E_API_DEMO_INVALID_NAME, "invalid name", query.CheckStringNotEmpty)
	qs.IntVar(&ap.Status, "status", false, errno.E_API_DEMO_INVALID_STATUS, "invalid status", nil)
	e := qs.Parse(context.QueryValues)

	if e != nil {
		return ap, nil, e
	}

	if ap.Status < 0 {
		return ap, nil, exception.New(errno.E_API_DEMO_INVALID_STATUS, "invalid status")
	}

	return ap, qs.ExistsInfo(), nil
}
