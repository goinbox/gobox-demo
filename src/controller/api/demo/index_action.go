package demo

import (
	"gdemo/controller/api"
	"gdemo/misc"
	"gdemo/model/demo"
	"gdemo/pcontext"
	"gdemo/perror"
	"gdemo/tasks/api/demo/list"
)

type indexRequest struct {
	IDs    []int64 `validate:"omitempty,min=1,max=100,dive,min=1"`
	Status *int    `validate:"omitempty,demo_status"`

	*misc.CommonListParams
}

type IndexResponse struct {
	Total int64

	DemoList []*demo.Entity
}

type indexAction struct {
	api.ApiAction

	req  *indexRequest
	resp *IndexResponse
}

func newIndexAction() *indexAction {
	a := &indexAction{
		req: &indexRequest{
			CommonListParams: misc.NewDefaultCommonListParams(),
		},
		resp: &IndexResponse{
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

func (a *indexAction) Run(ctx *pcontext.Context) error {
	out := &list.TaskOut{}
	err := api.RunTask(ctx, list.NewTask(), &list.TaskIn{
		IDs:        a.req.IDs,
		Status:     a.req.Status,
		ListParams: a.req.CommonListParams,
	}, out)
	if err != nil {
		return perror.New(perror.ECommonSysError, "sys error")
	}

	a.resp.Total = out.Total
	if out.DemoList != nil {
		a.resp.DemoList = out.DemoList
	}

	return nil
}
