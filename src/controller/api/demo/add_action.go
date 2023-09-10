package demo

import (
	"gdemo/controller/api"
	"gdemo/model"
	"gdemo/model/demo"
	"gdemo/pcontext"
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
	api.ApiAction

	req  *addRequest
	resp *AddResponse
}

func newAddAction() *addAction {
	a := &addAction{
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

func (a *addAction) Run(ctx *pcontext.Context) error {
	out := &add.TaskOut{}
	err := api.RunTask(ctx, add.NewTask(), &add.TaskIn{
		Name:   a.req.Name,
		Status: demo.StatusOnline,
	}, out)
	if err != nil {
		if model.DuplicateError(err) {
			return perror.New(perror.ECommonDataAlreadyExist, "data already exist")
		} else {
			return perror.New(perror.ECommonSysError, "sys error")
		}
	}

	a.resp.ID = out.ID
	return nil
}
