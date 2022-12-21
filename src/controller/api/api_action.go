package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/goinbox/golog"

	"gdemo/conf"
	"gdemo/controller"
	"gdemo/logic/factory"
	"gdemo/misc"
	"gdemo/perror"
	"gdemo/validate"
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
	if a.needSign() {
		err := a.checkSign()
		if err != nil {
			JumpToErrorAction(a.BaseAction, perror.NewFromError(perror.ECommonAuthFailure,
				fmt.Errorf("checkSign error: %w", err)))
		}
	}

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

func (a *ApiAction) needSign() bool {
	if !conf.ServerConf.Api.EnableSign {
		return false
	}

	if a.RemoteRealAddr.IP == "127.0.0.1" {
		return false
	}

	return true
}

type signParams struct {
	token string
	ts    int64
	sign  string
}

func (a *ApiAction) checkSign() error {
	params, err := a.parseSignParams()
	if err != nil {
		return fmt.Errorf("parseSignParams error: %w", err)
	}

	diff := time.Now().Unix() - params.ts
	abs := int(math.Abs(float64(diff)))
	if abs > 300 {
		return fmt.Errorf("check Timestamp failed")
	}

	sign := misc.ApiSign(params.token, params.ts, a.ReqRawBody)
	if sign != params.sign {
		return fmt.Errorf("check Sign failed, client sign is %s, server sign is %s", params.sign, sign)
	}

	return nil
}

func (a *ApiAction) parseSignParams() (*signParams, error) {
	header := a.Request().Header
	appName := header.Get("App")
	if appName == "" {
		return nil, fmt.Errorf("miss header App")
	}
	app := factory.DefaultLogicFactory.AppLogic().AppByName(a.Ctx, appName)
	if app == nil {
		return nil, fmt.Errorf("app %s not exist", appName)
	}

	tsStr := header.Get("Timestamp")
	if tsStr == "" {
		return nil, fmt.Errorf("miss header Timestamp")
	}
	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Timestamp.strconv.ParseInt error: %w", err)
	}

	sign := header.Get("Sign")
	if sign == "" {
		return nil, fmt.Errorf("miss header Sign")
	}

	a.Ctx.Logger.Info("signParams", []*golog.Field{
		{
			Key:   "token",
			Value: app.Token,
		},
		{
			Key:   "ts",
			Value: ts,
		},
		{
			Key:   "sign",
			Value: sign,
		},
	}...)

	return &signParams{
		token: app.Token,
		ts:    ts,
		sign:  sign,
	}, nil
}
