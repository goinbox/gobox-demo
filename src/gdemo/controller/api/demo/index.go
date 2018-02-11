package demo

import (
	"gdemo/dao"
	"gdemo/errno"
	"gdemo/svc"

	"github.com/goinbox/exception"
	"github.com/goinbox/gohttp/query"
)

type indexActionParams struct {
	Status int `mysql:"status"`

	offset int64
	cnt    int64
}

var indexQueryConditions map[string]string = map[string]string{
	"status": dao.SQL_COND_EQUAL,
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
		context.ApiData.Err = exception.New(errno.E_SYS_MYSQL_ERROR, err.Error())
		return
	}

	context.ApiData.Data = entities
}

func (d *DemoController) parseIndexActionParams(context *DemoContext) (*indexActionParams, map[string]bool, *exception.Exception) {
	ap := new(indexActionParams)

	qs := query.NewQuerySet()
	qs.IntVar(&ap.Status, "status", false, errno.E_API_DEMO_INVALID_STATUS, "invalid status", nil)
	qs.Int64Var(&ap.offset, "offset", false, errno.E_COMMON_INVALID_QUERY_OFFSET, "invalid offset", nil)
	qs.Int64Var(&ap.cnt, "cnt", false, errno.E_COMMON_INVALID_QUERY_CNT, "invalid cnt", nil)
	e := qs.Parse(context.QueryValues)
	if e != nil {
		return ap, nil, e
	}

	if ap.Status < 0 {
		return ap, nil, exception.New(errno.E_API_DEMO_INVALID_STATUS, "invalid status")
	}
	if ap.offset < 0 {
		return ap, nil, exception.New(errno.E_COMMON_INVALID_QUERY_OFFSET, "invalid offset")
	}
	if ap.cnt < 0 {
		return ap, nil, exception.New(errno.E_COMMON_INVALID_QUERY_CNT, "invalid cnt")
	}

	return ap, qs.ExistsInfo(), nil
}
