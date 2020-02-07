package mongodemo

import (
	"gdemo/errno"
	demoSvc "gdemo/svc/mongodemo"

	"github.com/goinbox/exception"
	"github.com/goinbox/gohttp/query"
)

func (d *MongoDemoController) AddAction(context *MongoDemoContext) {
	ap, e := d.parseAddActionParams(context)
	if e != nil {
		context.ApiData.Err = e
		return
	}

	ids, err := context.demoSvc.Insert(ap)
	if err != nil {
		context.ApiData.Err = exception.New(errno.E_COMMON_INSERT_ENTITY_FAILED, err.Error())
		return
	}

	context.ApiData.Data = map[string]interface{}{
		"id": ids[0],
	}
}

func (d *MongoDemoController) parseAddActionParams(context *MongoDemoContext) (*demoSvc.MongoDemoEntity, *exception.Exception) {
	ap := new(demoSvc.MongoDemoEntity)

	qs := query.NewQuerySet()
	qs.StringVar(&ap.Name, "name", true, errno.E_COMMON_INVALID_ARG, "invalid name", query.CheckStringNotEmpty)
	qs.IntVar(&ap.Status, "status", true, errno.E_COMMON_INVALID_ARG, "invalid status", nil)
	e := qs.Parse(context.QueryValues)
	if e != nil {
		return ap, e
	}

	if ap.Status < 0 {
		return ap, exception.New(errno.E_COMMON_INVALID_ARG, "invalid status")
	}

	return ap, nil
}
