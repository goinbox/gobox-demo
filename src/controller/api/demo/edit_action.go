package demo

import (
	"net/http"

	"github.com/goinbox/golog"

	"gdemo/controller/api"
	"gdemo/logic/factory"
	"gdemo/model"
	"gdemo/model/demo"
	"gdemo/perror"
)

type editRequest struct {
	ID     int64   `validate:"required,min=1"`
	Name   *string `validate:"omitempty,min=1"`
	Status *int    `validate:"omitempty,demo_status"`
}

type editAction struct {
	*api.ApiAction

	req *editRequest
}

func newEditAction(r *http.Request, w http.ResponseWriter, args []string) *editAction {
	a := &editAction{
		ApiAction: api.NewApiAction(r, w, args),

		req: new(editRequest),
	}

	a.RequestData = a.req

	return a
}

func (a *editAction) Name() string {
	return "edit"
}

func (a *editAction) Run() {
	fields := make(map[string]interface{})
	if a.req.Name != nil {
		fields[demo.ColumnName] = *a.req.Name
	}
	if a.req.Status != nil {
		fields[demo.ColumnStatus] = *a.req.Status
	}
	if len(fields) == 0 {
		return
	}

	logic := factory.DefaultLogicFactory.DemoLogic()
	err := logic.UpdateByIDs(a.Ctx, fields, a.req.ID)
	if err == nil {
		return
	}

	a.Ctx.Logger.Error("logic.UpdateByIDs error", golog.ErrorField(err))

	if model.DuplicateError(err) {
		a.Err = perror.New(perror.ECommonDataAlreadyExist, "data already exist")
	} else {
		a.Err = perror.New(perror.ECommonSysError, "sys error")
	}
}
