package demo

import (
	"gdemo/controller/api"
	demoSvc "gdemo/svc/demo"

	gcontroller "github.com/goinbox/gohttp/controller"

	"net/http"
)

type DemoContext struct {
	*api.ApiContext

	demoSvc *demoSvc.DemoSvc
}

func (d *DemoContext) BeforeAction() {
	d.ApiContext.BeforeAction()

	d.demoSvc = demoSvc.NewDemoSvc(d.TraceId, d.AccessLogger)
}

type DemoController struct {
	api.BaseController
}

func (d *DemoController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) gcontroller.ActionContext {
	context := new(DemoContext)
	context.ApiContext = d.BaseController.NewActionContext(req, respWriter).(*api.ApiContext)

	return context
}
