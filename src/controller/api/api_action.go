package api

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/goinbox/golog"

	"gdemo/conf"
	"gdemo/controller"
	"gdemo/logic/factory"
	"gdemo/misc"
	"gdemo/pcontext"
	"gdemo/perror"
	"gdemo/validate"
)

type BaseResponse struct {
	Errno int
	Msg   string
	Tid   string
}

type Response struct {
	*BaseResponse

	Data interface{} `json:",omitempty"`
}

type ApiAction struct {
	controller.BaseAction

	RequestData  interface{}
	ResponseData interface{}
}

func (a *ApiAction) Before(ctx *pcontext.Context) error {
	if a.needSign() {
		err := a.checkSign(ctx)
		if err != nil {
			return perror.NewFromError(perror.ECommonAuthFailure,
				fmt.Errorf("checkSign error: %w", err))
		}
	}

	err := json.Unmarshal(a.ReqRawBody, a.RequestData)
	if err != nil {
		return perror.NewFromError(perror.ECommonInvalidArg,
			fmt.Errorf("json.Unmarshal ReqRawBody error: %w", err))
	}

	err = validate.Validator().Struct(a.RequestData)
	if err != nil {
		return perror.NewFromError(perror.ECommonInvalidArg,
			fmt.Errorf("check args error: %w", err))
	}

	return nil
}

func (a *ApiAction) After(ctx *pcontext.Context, err error) {
	resp := &Response{
		BaseResponse: &BaseResponse{
			Errno: perror.Success,
			Msg:   "",
			Tid:   ctx.TraceID(),
		},
	}

	if err != nil {
		perr := perror.ParsePerror(err)
		if perr != nil {
			resp.Errno = perr.Errno()
			resp.Msg = perr.Msg()
		} else {
			resp.Errno = perror.ECommonSysError
			resp.Msg = err.Error()
		}
	} else {
		resp.Data = a.ResponseData
	}

	body, _ := json.Marshal(resp)
	a.SetResponseBody(body)

	a.BaseAction.After(ctx, err)
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

func (a *ApiAction) checkSign(ctx *pcontext.Context) error {
	params, err := a.parseSignParams(ctx)
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

func (a *ApiAction) parseSignParams(ctx *pcontext.Context) (*signParams, error) {
	header := a.Request().Header
	appName := header.Get("App")
	if appName == "" {
		return nil, fmt.Errorf("miss header App")
	}
	app := factory.DefaultLogicFactory.AppLogic().AppByName(ctx, appName)
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

	ctx.Logger().Info("signParams", []*golog.Field{
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
