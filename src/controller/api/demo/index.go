package demo

import (
	"gdemo/errno"
	"gdemo/svc"

	"github.com/goinbox/goerror"
	"github.com/goinbox/gohttp/query"
	"github.com/goinbox/mysql"
)

type indexActionParams struct {
	Status int `mysql:"status"`

	offset int64
	cnt    int64
}

var indexQueryConditions map[string]string = map[string]string{
	"status": mysql.SqlCondEqual,
}

func (d *DemoController) IndexAction(context *DemoContext) {
	ap, exists, e := d.parseIndexActionParams(context)
	if e != nil {
		context.ApiData.Err = e
		return
	}

	sqp := &svc.SqlQueryParams{
		ParamsStructPtr: ap,
		Exists:          exists,
		Conditions:      indexQueryConditions,

		OrderBy: "id desc",
		Offset:  ap.offset,
		Cnt:     ap.cnt,
	}

	entities, err := context.demoSvc.SimpleQueryAnd(sqp)
	if err != nil {
		context.ApiData.Err = goerror.New(errno.ESysMysqlError, err.Error())
		return
	}

	context.ApiData.Data = entities
}

func (d *DemoController) parseIndexActionParams(context *DemoContext) (*indexActionParams, map[string]bool, *goerror.Error) {
	ap := new(indexActionParams)

	qs := query.NewQuerySet()
	qs.IntVar(&ap.Status, "status", false, errno.ECommonInvalidArg, "invalid status", nil)
	qs.Int64Var(&ap.offset, "offset", false, errno.ECommonInvalidArg, "invalid offset", nil)
	qs.Int64Var(&ap.cnt, "cnt", false, errno.ECommonInvalidArg, "invalid cnt", nil)
	e := qs.Parse(context.QueryValues)
	if e != nil {
		return ap, nil, e
	}

	if ap.Status < 0 {
		return ap, nil, goerror.New(errno.ECommonInvalidArg, "invalid status")
	}
	if ap.offset < 0 {
		return ap, nil, goerror.New(errno.ECommonInvalidArg, "invalid offset")
	}
	if ap.cnt < 0 {
		return ap, nil, goerror.New(errno.ECommonInvalidArg, "invalid cnt")
	}

	return ap, qs.ExistsInfo(), nil
}
