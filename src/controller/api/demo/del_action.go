package demo

import (
	"net/http"

	"gdemo/controller/api"
	"gdemo/perror"
	"gdemo/tasks/api/demo/del"
)

type delRequest struct {
	IDs []int64 `validate:"required,min=1,dive,min=1"`
}

type DelResponse struct {
	RowsAffected int64
}

type delAction struct {
	*api.ApiAction

	req  *delRequest
	resp *DelResponse
}

func newDelAction(r *http.Request, w http.ResponseWriter, args []string) *delAction {
	a := &delAction{
		ApiAction: api.NewApiAction(r, w, args),

		req:  new(delRequest),
		resp: new(DelResponse),
	}

	a.RequestData = a.req
	a.ResponseData = a.resp

	return a
}

func (a *delAction) Name() string {
	return "Del"
}

func (a *delAction) Run() {
	out := &del.TaskOut{}
	err := api.RunTask(a.Ctx, del.NewTask(), &del.TaskIn{
		IDs: a.req.IDs,
	}, out)
	if err == nil {
		a.resp.RowsAffected = out.RowsAffected
		return
	}

	a.Err = perror.New(perror.ECommonSysError, "sys error")
}
