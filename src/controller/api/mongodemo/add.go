package mongodemo

import (
	"github.com/goinbox/goerror"
	"github.com/goinbox/gohttp/query"

	"gdemo/define/entity"
	"gdemo/errno"
)

func (d *MongoDemoController) AddAction(context *MongoDemoContext) {
	ap, e := d.parseAddActionParams(context)
	if e != nil {
		context.ApiData.Err = e
		return
	}

	ids, err := context.demoSvc.Insert(ap)
	if err != nil {
		context.ApiData.Err = goerror.New(errno.ECommonInsertEntityFailed, err.Error())
		return
	}

	context.ApiData.Data = map[string]interface{}{
		"id": ids[0],
	}
}

func (d *MongoDemoController) parseAddActionParams(context *MongoDemoContext) (*entity.MongoDemoEntity, *goerror.Error) {
	ap := new(entity.MongoDemoEntity)

	qs := query.NewQuerySet()
	qs.StringVar(&ap.Name, "name", true, errno.ECommonInvalidArg, "invalid name", query.CheckStringNotEmpty)
	qs.IntVar(&ap.Status, "status", true, errno.ECommonInvalidArg, "invalid status", nil)
	e := qs.Parse(context.QueryValues)
	if e != nil {
		return ap, e
	}

	if ap.Status < 0 {
		return ap, goerror.New(errno.ECommonInvalidArg, "invalid status")
	}

	return ap, nil
}
