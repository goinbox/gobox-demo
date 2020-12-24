package mongodemo

import (
	"github.com/goinbox/goerror"
	"github.com/goinbox/gohttp/query"
	"github.com/goinbox/mongo"

	"gdemo/define"
	"gdemo/errno"
)

type indexActionParams struct {
	Status int `bson:"status"`

	offset int
	cnt    int
}

var indexQueryConditions map[string]string = map[string]string{
	"status": mongo.MONGO_COND_GREATER_EQUAL,
}

func (d *MongoDemoController) IndexAction(context *MongoDemoContext) {
	ap, exists, e := d.parseIndexActionParams(context)
	if e != nil {
		context.ApiData.Err = e
		return
	}

	mqp := &define.MongoQueryParams{
		ParamsStructPtr: ap,
		Exists:          exists,
		Conditions:      indexQueryConditions,

		OrderBy: []string{"-_id"},
		Offset:  ap.offset,
		Cnt:     ap.cnt,
	}

	entities, err := context.demoSvc.SelectAll(mqp)
	if err != nil {
		context.ApiData.Err = goerror.New(errno.ESysMongoError, err.Error())
		return
	}

	context.ApiData.Data = entities
}

func (d *MongoDemoController) parseIndexActionParams(context *MongoDemoContext) (*indexActionParams, map[string]bool, *goerror.Error) {
	ap := new(indexActionParams)

	qs := query.NewQuerySet()
	qs.IntVar(&ap.Status, "status", false, errno.ECommonInvalidArg, "invalid status", nil)
	qs.IntVar(&ap.offset, "offset", false, errno.ECommonInvalidArg, "invalid offset", nil)
	qs.IntVar(&ap.cnt, "cnt", false, errno.ECommonInvalidArg, "invalid cnt", nil)
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
