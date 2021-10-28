package api

import (
	"html"
	"net/http"

	"gdemo/perror"
	gcontroller "github.com/goinbox/gohttp/controller"

	"gdemo/controller"
	"gdemo/misc"
)

type IApiDataContext interface {
	gcontroller.ActionContext

	Version() string
	Data() interface{}
	Err() *perror.Error
}

type ApiContext struct {
	*controller.BaseContext

	ApiData struct {
		V    string
		Data interface{}
		Err  *perror.Error
	}
}

func (a *ApiContext) Version() string {
	return a.ApiData.V
}

func (a *ApiContext) Data() interface{} {
	return a.ApiData.Data
}

func (a *ApiContext) Err() *perror.Error {
	return a.ApiData.Err
}

func (a *ApiContext) AfterAction() {
	f := a.QueryValues.Get("fmt")
	if f == "jsonp" {
		callback := a.QueryValues.Get("_callback")
		if callback != "" {
			a.SetResponseBody(misc.ApiJsonp(a.ApiData.V, a.ApiData.Data, a.ApiData.Err, html.EscapeString(callback)))
			return
		}
	}

	a.SetResponseBody(misc.ApiJson(a.ApiData.V, a.ApiData.Data, a.ApiData.Err))
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
