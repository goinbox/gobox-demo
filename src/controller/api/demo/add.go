package demo

import (
	"gdemo/perror"
	"github.com/goinbox/gohttp/query"

	"gdemo/define/entity"
	"gdemo/perror"
)

func (d *DemoController) AddAction(context *DemoContext) {
	ap, e := d.parseAddActionParams(context)
	if e != nil {
		context.ApiData.Err = e
		return
	}

	ids, err := context.demoSvc.Insert(ap)
	if err != nil {
		context.ApiData.Err = perror.Error(perror.ECommonInsertEntityFailed, err.Error())
		return
	}

	context.ApiData.Data = map[string]interface{}{
		"id": ids[0],
	}
}

func (d *DemoController) parseAddActionParams(context *DemoContext) (*entity.DemoEntity, *perror.Error) {
	ap := new(entity.DemoEntity)

	qs := query.NewQuerySet()
	qs.StringVar(&ap.Name, "name", true, perror.ECommonInvalidArg, "invalid name", query.CheckStringNotEmpty)
	qs.IntVar(&ap.Status, "status", true, perror.ECommonInvalidArg, "invalid status", nil)
	e := qs.Parse(context.QueryValues)
	if e != nil {
		return ap, e
	}

	if ap.Status < 0 {
		return ap, perror.Error(perror.ECommonInvalidArg, "invalid status")
	}

	return ap, nil
}
