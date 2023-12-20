package demo

import (
	"gdemo/controller/api"
	"gdemo/model"
	"gdemo/pcontext"
	"gdemo/perror"
	"gdemo/tasks"
	"gdemo/tasks/api/demo/edit"
)

type editRequest struct {
	ID     int64   `validate:"required,min=1" mysql_update_column:"omit"`
	Name   *string `validate:"omitempty,min=1"`
	Status *int    `validate:"omitempty,demo_status"`
}

type EditResponse struct {
	RowsAffected int64
}

type editAction struct {
	api.ApiAction

	req  *editRequest
	resp *EditResponse
}

func newEditAction() *editAction {
	a := &editAction{
		req:  new(editRequest),
		resp: new(EditResponse),
	}

	a.RequestData = a.req
	a.ResponseData = a.resp

	return a
}

func (a *editAction) Name() string {
	return "Edit"
}

func (a *editAction) Run(ctx *pcontext.Context) error {
	out := &edit.TaskOut{}
	err := tasks.RunTask(ctx, edit.NewTask(), &edit.TaskIn{
		ID:           a.req.ID,
		UpdateParams: a.req,
	}, out)
	if err != nil {
		if model.DuplicateError(err) {
			return perror.New(perror.ECommonDataAlreadyExist, "data already exist")
		} else {
			return perror.New(perror.ECommonSysError, "sys error")
		}
	}

	a.resp.RowsAffected = out.RowsAffected
	return nil
}
