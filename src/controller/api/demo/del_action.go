package demo

import (
	"net/http"

	"github.com/goinbox/golog"

	"gdemo/controller/api"
	"gdemo/logic/factory"
	"gdemo/perror"
)

type delRequest struct {
	IDs []int64 `validate:"required,min=1,dive,min=1"`
}

type delAction struct {
	*api.ApiAction

	req *delRequest
}

func newDelAction(r *http.Request, w http.ResponseWriter, args []string) *delAction {
	a := &delAction{
		ApiAction: api.NewApiAction(r, w, args),

		req: new(delRequest),
	}

	a.RequestData = a.req

	return a
}

func (a *delAction) Name() string {
	return "del"
}

func (a *delAction) Run() {
	logic := factory.DefaultLogicFactory.DemoLogic()
	err := logic.DeleteByIDs(a.Ctx, a.req.IDs...)
	if err == nil {
		return
	}

	a.Ctx.Logger.Error("logic.DeleteByIDs error", golog.ErrorField(err))
	a.Err = perror.New(perror.ECommonSysError, "sys error")
}
