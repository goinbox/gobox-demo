package mongodemo

import (
	"github.com/goinbox/goerror"
	"github.com/goinbox/gohttp/query"

	"gdemo/define/entity"
	"gdemo/errno"
)

type editActionParams struct {
	item entity.MongoDemoEntity

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

func (d *MongoDemoController) parseEditActionParams(context *MongoDemoContext) (*entity.MongoDemoEntity, map[string]bool, *goerror.Error) {
	ap := new(editActionParams)

	qs := query.NewQuerySet()

	qs.Int64Var(&ap.tid, "id", true, errno.ECommonInvalidArg, "invalid id", query.CheckInt64IsPositive)
	qs.StringVar(&ap.item.Name, "name", false, errno.ECommonInvalidArg, "invalid name", query.CheckStringNotEmpty)
	qs.IntVar(&ap.item.Status, "status", false, errno.ECommonInvalidArg, "invalid status", nil)
	e := qs.Parse(context.QueryValues)

	if e != nil {
		return &ap.item, nil, e
	}

	ap.item.Id = ap.tid

	if ap.item.Status < 0 {
		return &ap.item, nil, goerror.New(errno.ECommonInvalidArg, "invalid status")
	}

	return &ap.item, qs.ExistsInfo(), nil
}
