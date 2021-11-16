package api

import (
	"html"
	"net/http"

	gcontroller "github.com/goinbox/gohttp/controller"

	"gdemo/perror"

	"gdemo/controller"
)

type ApiContext interface {
	gcontroller.ActionContext

	Data() interface{}
	Err() *perror.Error
}

type BaseApiContext struct {
	*controller.BaseContext

	ApiData struct {
		Data interface{}
		Err  *perror.Error
	}
}

func (a *BaseApiContext) Data() interface{} {
	return a.ApiData.Data
}

func (a *BaseApiContext) Err() *perror.Error {
	return a.ApiData.Err
}

func (a *BaseApiContext) AfterAction() {
	f := a.QueryValues.Get("fmt")
	if f == "jsonp" {
		callback := a.QueryValues.Get("_callback")
		if callback != "" {
			a.SetResponseBody(ApiJsonp(a.ApiData.Data, a.ApiData.Err, html.EscapeString(callback)))
			return
		}
	}

	a.SetResponseBody(ApiJson(a.ApiData.Data, a.ApiData.Err))
	a.BaseContext.AfterAction()
}

type ApiController struct {
	controller.BaseController
}

func (c *ApiController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) gcontroller.ActionContext {
	context := new(BaseApiContext)
	context.BaseContext = c.BaseController.NewActionContext(req, respWriter).(*controller.BaseContext)

	return context
}

func JumpToApiError(context gcontroller.ActionContext, args ...interface{}) {
	acontext := context.(ApiContext)

	acontext.SetResponseBody(ApiJson(acontext.Data(), acontext.Err()))
}
