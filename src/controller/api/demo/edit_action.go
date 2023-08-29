package demo

import (
	"net/http"

	"gdemo/controller/api"
	"gdemo/model"
	"gdemo/perror"
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
	*api.ApiAction

	req  *editRequest
	resp *EditResponse
}

func newEditAction(r *http.Request, w http.ResponseWriter, args []string) *editAction {
	a := &editAction{
		ApiAction: api.NewApiAction(r, w, args),

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

func (a *editAction) Run() {
	out := &edit.TaskOut{}
	err := api.RunTask(a.Ctx, edit.NewTask(), &edit.TaskIn{
		ID:           a.req.ID,
		UpdateParams: a.req,
	}, out)
	if err == nil {
		a.resp.RowsAffected = out.RowsAffected
		return
	}

	if model.DuplicateError(err) {
		a.Err = perror.New(perror.ECommonDataAlreadyExist, "data already exist")
	} else {
		a.Err = perror.New(perror.ECommonSysError, "sys error")
	}
}
