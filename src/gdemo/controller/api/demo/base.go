package demo

import (
	"gdemo/conf"
	"gdemo/controller/api"
	"gdemo/svc"

	gcontroller "github.com/goinbox/gohttp/controller"

	"net/http"
)

type DemoContext struct {
	*api.ApiContext

	demoSvc *svc.DemoSvc
}

func (this *DemoContext) BeforeAction() {
	this.ApiContext.BeforeAction()

	this.demoSvc = svc.NewDemoSvc(
		this.ErrorLogger,
		this.MysqlClient,
		conf.BaseConf.PrjName,
		this.RedisClient,
	)
}

type DemoController struct {
	api.BaseController
}

func (this *DemoController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) gcontroller.ActionContext {
	context := new(DemoContext)
	context.ApiContext = this.BaseController.NewActionContext(req, respWriter).(*api.ApiContext)

	return context
}
