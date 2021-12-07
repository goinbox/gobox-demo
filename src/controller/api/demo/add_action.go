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

type addRequest struct {
	Name string `validate:"required,min=1"`
}

type addAction struct {
	*api.ApiAction

	req *addRequest
}

func newAddAction(r *http.Request, w http.ResponseWriter, args []string) *addAction {
	a := &addAction{
		ApiAction: api.NewApiAction(r, w, args),

		req: new(addRequest),
	}

	a.RequestData = a.req

	return a
}

func (a *addAction) Name() string {
	return "add"
}

func (a *addAction) Run() {
	entity := &demo.Entity{
		BaseEntity: model.BaseEntity{},
		Name:       a.req.Name,
		Status:     demo.StatusOnline,
	}

	logic := factory.DefaultLogicFactory.DemoLogic()
	err := logic.Insert(a.Ctx, entity)
	if err == nil {
		return
	}

	a.Ctx.Logger.Error("logic.Insert error", golog.ErrorField(err))

	if model.DuplicateError(err) {
		a.Err = perror.New(perror.ECommonDataAlreadyExist, "data already exist")
	} else {
		a.Err = perror.New(perror.ECommonSysError, "sys error")
	}
}
