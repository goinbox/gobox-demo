package api

import (
	"html"
	"net/http"

	gcontroller "github.com/goinbox/gohttp/controller"

	"gdemo/perror"

	"gdemo/controller"
)

type ApiData struct {
	Err  *perror.Error
	Tid  string
	Data interface{}
}

type ApiContext struct {
	*controller.BaseContext

	ApiData *ApiData
}

func (a *ApiContext) AfterAction() {
	f := a.QueryValues.Get("fmt")
	if f == "jsonp" {
		callback := a.QueryValues.Get("_callback")
		if callback != "" {
			a.SetResponseBody(ApiJsonp(a.ApiData, html.EscapeString(callback)))
			return
		}
	}

	a.SetResponseBody(ApiJson(a.ApiData))
	a.BaseContext.AfterAction()
}

type ApiController struct {
	controller.BaseController
}

func (c *ApiController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) gcontroller.ActionContext {
	bctx := c.BaseController.NewActionContext(req, respWriter).(*controller.BaseContext)

	return &ApiContext{
		BaseContext: bctx,
		ApiData: &ApiData{
			Tid: bctx.Ctx.TraceID,
		},
	}
}
