package demo

import (
	"net/http"

	"gdemo/controller/api"
	"gdemo/model"
	"gdemo/model/demo"
	"gdemo/perror"
	"gdemo/tasks/api/demo/add"
)

type addRequest struct {
	Name string `validate:"required,min=1"`
}

type AddResponse struct {
	ID int64
}

type addAction struct {
	*api.ApiAction

	req  *addRequest
	resp *AddResponse
}

func newAddAction(r *http.Request, w http.ResponseWriter, args []string) *addAction {
	a := &addAction{
		ApiAction: api.NewApiAction(r, w, args),

		req:  new(addRequest),
		resp: new(AddResponse),
	}

	a.RequestData = a.req
	a.ResponseData = a.resp

	return a
}

func (a *addAction) Name() string {
	return "Add"
}

func (a *addAction) Run() {
	out := &add.TaskOut{}
	err := api.RunTask(a.Ctx, add.NewTask(a.Ctx), &add.TaskIn{
		Name:   a.req.Name,
		Status: demo.StatusOnline,
	}, out)
	if err == nil {
		a.resp.ID = out.ID
		return
	}

	if model.DuplicateError(err) {
		a.Err = perror.New(perror.ECommonDataAlreadyExist, "data already exist")
	} else {
		a.Err = perror.New(perror.ECommonSysError, "sys error")
	}
}
