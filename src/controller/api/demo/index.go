package demo

import (
	"github.com/goinbox/goerror"
	"github.com/goinbox/gohttp/query"
	"github.com/goinbox/mysql"

	"gdemo/define"
	"gdemo/perror"
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
	ap, _, e := d.parseIndexActionParams(context)
	if e != nil {
		context.ApiData.Err = e
		return
	}

	sqp := &define.SqlQueryParams{
		CondItems: []*mysql.SqlColQueryItem{
			{
				Name:      "status",
				Condition: mysql.SqlCondEqual,
				Value:     ap.Status,
			},
		},

		OrderBy: "id desc",
		Offset:  ap.offset,
		Cnt:     ap.cnt,
	}

	entities, err := context.demoSvc.SimpleQueryAnd(sqp)
	if err != nil {
		context.ApiData.Err = goerror.New(perror.ESysMysqlError, err.Error())
		return
	}

	context.ApiData.Data = entities
}

func (d *DemoController) parseIndexActionParams(context *DemoContext) (*indexActionParams, map[string]bool, *goerror.Error) {
	ap := new(indexActionParams)

	qs := query.NewQuerySet()
	qs.IntVar(&ap.Status, "status", false, perror.ECommonInvalidArg, "invalid status", nil)
	qs.Int64Var(&ap.offset, "offset", false, perror.ECommonInvalidArg, "invalid offset", nil)
	qs.Int64Var(&ap.cnt, "cnt", false, perror.ECommonInvalidArg, "invalid cnt", nil)
	e := qs.Parse(context.QueryValues)
	if e != nil {
		return ap, nil, e
	}

	if ap.Status < 0 {
		return ap, nil, goerror.New(perror.ECommonInvalidArg, "invalid status")
	}
	if ap.offset < 0 {
		return ap, nil, goerror.New(perror.ECommonInvalidArg, "invalid offset")
	}
	if ap.cnt < 0 {
		return ap, nil, goerror.New(perror.ECommonInvalidArg, "invalid cnt")
	}

	return ap, qs.ExistsInfo(), nil
}
