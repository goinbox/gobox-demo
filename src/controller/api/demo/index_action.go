package demo

import (
	"net/http"

	"github.com/goinbox/golog"
	"github.com/goinbox/mysql"

	"gdemo/controller/api"
	"gdemo/logic/factory"
	"gdemo/model"
	"gdemo/model/demo"
	"gdemo/perror"
)

type indexRequest struct {
	IDs    []int64 `validate:"omitempty,min=1,max=100,dive,min=1"`
	Status *int    `validate:"omitempty,demo_status"`

	Offset int64
	Cnt    int64
}

type indexResponse struct {
	Total int64

	DemoList []*demo.Entity
}

type indexAction struct {
	*api.ApiAction

	req  *indexRequest
	resp *indexResponse
}

func newIndexAction(r *http.Request, w http.ResponseWriter, args []string) *indexAction {
	a := &indexAction{
		ApiAction: api.NewApiAction(r, w, args),

		req: &indexRequest{
			Offset: 0,
			Cnt:    10,
		},
		resp: new(indexResponse),
	}

	a.RequestData = a.req
	a.ResponseData = a.resp

	return a
}

func (a *indexAction) Name() string {
	return "Index"
}

func (a *indexAction) Run() {
	logic := factory.DefaultLogicFactory.DemoLogic()

	conds := a.sqlConds()

	total, err := logic.SimpleTotalAnd(a.Ctx, conds...)
	if err != nil {
		a.Ctx.Logger.Error("logic.SimpleTotalAnd error", golog.ErrorField(err))
		a.Err = perror.New(perror.ECommonSysError, "sys error")
		return
	}

	params := &mysql.SqlQueryParams{
		CondItems: conds,

		OrderBy: "id desc",
		Offset:  a.req.Offset,
		Cnt:     a.req.Cnt,
	}
	entities, err := logic.SimpleQueryAnd(a.Ctx, params)
	if err != nil {
		a.Ctx.Logger.Error("logic.SimpleQueryAnd error", golog.ErrorField(err))
		a.Err = perror.New(perror.ECommonSysError, "sys error")
		return
	}

	a.resp.Total = total
	a.resp.DemoList = entities
}

func (a *indexAction) sqlConds() []*mysql.SqlColQueryItem {
	var conds []*mysql.SqlColQueryItem

	item := model.MakeInt64SliceSqlColQueryItem(demo.ColumnID, a.req.IDs)
	if item != nil {
		conds = append(conds, item)
	}
	if a.req.Status != nil {
		conds = append(conds, &mysql.SqlColQueryItem{
			Name:      demo.ColumnStatus,
			Condition: mysql.SqlCondEqual,
			Value:     *a.req.Status,
		})
	}

	return conds
}
