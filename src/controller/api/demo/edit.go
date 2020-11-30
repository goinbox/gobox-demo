package demo

import (
	"gdemo/errno"
	demoSvc "gdemo/svc/demo"

	"github.com/goinbox/goerror"
	"github.com/goinbox/gohttp/query"
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

	updated, err := context.demoSvc.UpdateById(ap.Id, ap, exists)
	if err != nil {
		context.ApiData.Err = goerror.New(errno.ECommonUpdateEntityFailed, err.Error())
		return
	}

	context.ApiData.Data = updated
}

func (d *DemoController) parseEditActionParams(context *DemoContext) (*demoSvc.DemoEntity, map[string]bool, *goerror.Error) {
	ap := new(demoSvc.DemoEntity)

	qs := query.NewQuerySet()
	qs.Int64Var(&ap.Id, "id", true, errno.ECommonInvalidArg, "invalid id", query.CheckInt64IsPositive)
	qs.StringVar(&ap.Name, "name", false, errno.ECommonInvalidArg, "invalid name", query.CheckStringNotEmpty)
	qs.IntVar(&ap.Status, "status", false, errno.ECommonInvalidArg, "invalid status", nil)
	e := qs.Parse(context.QueryValues)
	if e != nil {
		return ap, nil, e
	}

	if ap.Status < 0 {
		return ap, nil, goerror.New(errno.ECommonInvalidArg, "invalid status")
	}

	return ap, qs.ExistsInfo(), nil
}
