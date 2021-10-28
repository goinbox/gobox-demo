package demo

import (
	"gdemo/perror"
	"github.com/goinbox/gohttp/query"

	"gdemo/define/entity"
	"gdemo/perror"
)

func (d *DemoController) EditAction(context *DemoContext) {
	ap, exists, e := d.parseEditActionParams(context)
	if e != nil {
		context.ApiData.Err = e
		return
	}

	if len(exists) == 1 { //only has id
		return
	}

	updateFields := make(map[string]interface{})
	_, ok := exists["name"]
	if ok {
		updateFields["name"] = ap.Name
	}
	_, ok = exists["status"]
	if ok {
		updateFields["status"] = ap.Status
	}

	updated, err := context.demoSvc.UpdateById(ap.Id, updateFields)
	if err != nil {
		context.ApiData.Err = perror.Error(perror.ECommonUpdateEntityFailed, err.Error())
		return
	}

	context.ApiData.Data = updated
}

func (d *DemoController) parseEditActionParams(context *DemoContext) (*entity.DemoEntity, map[string]bool, *perror.Error) {
	ap := new(entity.DemoEntity)

	qs := query.NewQuerySet()
	qs.Int64Var(&ap.Id, "id", true, perror.ECommonInvalidArg, "invalid id", query.CheckInt64IsPositive)
	qs.StringVar(&ap.Name, "name", false, perror.ECommonInvalidArg, "invalid name", query.CheckStringNotEmpty)
	qs.IntVar(&ap.Status, "status", false, perror.ECommonInvalidArg, "invalid status", nil)
	e := qs.Parse(context.QueryValues)
	if e != nil {
		return ap, nil, e
	}

	if ap.Status < 0 {
		return ap, nil, perror.Error(perror.ECommonInvalidArg, "invalid status")
	}

	return ap, qs.ExistsInfo(), nil
}
