package demo

import (
	"net/http"

	"github.com/goinbox/golog"

	"gdemo/controller/api"
	"gdemo/logic/factory"
	"gdemo/misc"
	"gdemo/model"
	"gdemo/perror"
)

type editRequest struct {
	ID     int64   `validate:"required,min=1" mysql_update_column:"omit"`
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
	return "Edit"
}

func (a *editAction) Run() {
	updateColumns := misc.MakeMysqlUpdateColumns(a.req)
	if len(updateColumns) == 0 {
		return
	}

	logic := factory.DefaultLogicFactory.DemoLogic()
	err := logic.UpdateByIDs(a.Ctx, updateColumns, a.req.ID)
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
