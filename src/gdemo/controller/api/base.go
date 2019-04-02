package api

import (
	"gdemo/controller"
	"gdemo/misc"

	"github.com/goinbox/exception"
	gcontroller "github.com/goinbox/gohttp/controller"

	"html"
	"net/http"
)

type IApiDataContext interface {
	gcontroller.ActionContext

	Version() string
	Data() interface{}
	Err() *exception.Exception
}

type ApiContext struct {
	*controller.BaseContext

	ApiData struct {
		V    string
		Data interface{}
		Err  *exception.Exception
	}
}

func (a *ApiContext) Version() string {
	return a.ApiData.V
}

func (a *ApiContext) Data() interface{} {
	return a.ApiData.Data
}

func (a *ApiContext) Err() *exception.Exception {
	return a.ApiData.Err
}

func (a *ApiContext) AfterAction() {
	f := a.QueryValues.Get("fmt")
	if f == "jsonp" {
		callback := a.QueryValues.Get("_callback")
		if callback != "" {
			a.RespBody = misc.ApiJsonp(a.ApiData.V, a.ApiData.Data, a.ApiData.Err, html.EscapeString(callback))
			return
		}
	}

	a.RespBody = misc.ApiJson(a.ApiData.V, a.ApiData.Data, a.ApiData.Err)
	a.BaseContext.AfterAction()
}

type BaseController struct {
	controller.BaseController
}

func (b *BaseController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) gcontroller.ActionContext {
	context := new(ApiContext)
	context.BaseContext = b.BaseController.NewActionContext(req, respWriter).(*controller.BaseContext)

	return context
}

func JumpToApiError(context gcontroller.ActionContext, args ...interface{}) {
	acontext := context.(IApiDataContext)

	acontext.SetResponseBody(misc.ApiJson(acontext.Version(), acontext.Data(), acontext.Err()))
}
