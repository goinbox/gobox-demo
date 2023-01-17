package api

import (
	"encoding/json"

	"gdemo/controller"
	"gdemo/perror"
)

type ErrorAction struct {
	*controller.BaseAction

	Err *perror.Error
}

func (a *ErrorAction) Name() string {
	return "Error"
}

func (a *ErrorAction) Run() {
	resp := &Response{
		BaseResponse: &BaseResponse{
			Errno: a.Err.Errno(),
			Msg:   a.Err.Msg(),
			Tid:   a.Ctx.TraceID,
		},
	}

	body, _ := json.Marshal(resp)
	a.SetResponseBody(body)
}

func JumpToErrorAction(a *controller.BaseAction, err *perror.Error) {
	panic(&ErrorAction{
		BaseAction: a,
		Err:        err,
	})
}
