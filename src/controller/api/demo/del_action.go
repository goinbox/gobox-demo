package demo

import (
	"gdemo/controller/api"
	"gdemo/pcontext"
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
	api.ApiAction

	req  *delRequest
	resp *DelResponse
}

func newDelAction() *delAction {
	a := &delAction{
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

func (a *delAction) Run(ctx *pcontext.Context) error {
	out := &del.TaskOut{}
	err := api.RunTask(ctx, del.NewTask(), &del.TaskIn{
		IDs: a.req.IDs,
	}, out)
	if err != nil {
		return perror.New(perror.ECommonSysError, "sys error")
	}

	a.resp.RowsAffected = out.RowsAffected
	return nil
}
