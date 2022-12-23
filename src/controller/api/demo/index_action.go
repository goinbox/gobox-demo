package demo

import (
	"net/http"

	"gdemo/controller/api"
	"gdemo/misc"
	"gdemo/model/demo"
	"gdemo/perror"
	"gdemo/task/api/demo/list"
)

type indexRequest struct {
	IDs    []int64 `validate:"omitempty,min=1,max=100,dive,min=1"`
	Status *int    `validate:"omitempty,demo_status"`

	*misc.CommonListParams
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
			CommonListParams: misc.NewDefaultCommonListParams(),
		},
		resp: &indexResponse{
			DemoList: []*demo.Entity{},
		},
	}

	a.RequestData = a.req
	a.ResponseData = a.resp

	return a
}

func (a *indexAction) Name() string {
	return "Index"
}

func (a *indexAction) Run() {
	out := &list.TaskOut{}
	err := api.RunTask(a.Ctx, list.NewTask(a.Ctx), &list.TaskIn{
		IDs:        a.req.IDs,
		Status:     a.req.Status,
		ListParams: a.req.CommonListParams,
	}, out)
	if err == nil {
		a.resp.Total = out.Total
		if out.DemoList != nil {
			a.resp.DemoList = out.DemoList
		}
	} else {
		a.Err = perror.New(perror.ECommonSysError, "sys error")
	}
}
