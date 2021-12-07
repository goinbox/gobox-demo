package api

import (
	"encoding/json"
	"net/http"

	"github.com/goinbox/golog"

	"gdemo/controller"
	"gdemo/controller/api/validate"
	"gdemo/perror"
)

type Response struct {
	Errno int
	Msg   string
	Tid   string

	Data interface{} `json:",omitempty"`
}

type ApiAction struct {
	*controller.BaseAction

	Err *perror.Error

	RequestData  interface{}
	ResponseData interface{}
}

func NewApiAction(r *http.Request, w http.ResponseWriter, args []string) *ApiAction {
	return &ApiAction{
		BaseAction: controller.NewBaseAction(r, w, args),
	}
}

func (a *ApiAction) Before() {
	err := json.Unmarshal(a.ReqRawBody, a.RequestData)
	if err != nil {
		a.Ctx.Logger.Warning("json.Unmarshal ReqRawBody error", golog.ErrorField(err))
	}

	err = validate.Validator().Struct(a.RequestData)
	if err != nil {
		JumpToErrorAction(a.BaseAction, perror.NewFromError(perror.ECommonInvalidArg, err))
	}
}

func (a *ApiAction) After() {
	resp := &Response{
		Errno: perror.Success,
		Msg:   "",
		Tid:   a.Ctx.TraceID,
	}

	if a.Err != nil {
		resp.Errno = a.Err.Errno()
		resp.Msg = a.Err.Msg()
	} else {
		resp.Data = a.ResponseData
	}

	body, _ := json.Marshal(resp)
	a.SetResponseBody(body)

	a.BaseAction.After()
}
