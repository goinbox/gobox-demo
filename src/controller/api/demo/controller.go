package demo

import (
	"gdemo/controller/api"
	"gdemo/logic/demo"
	"gdemo/logic/factory"

	gcontroller "github.com/goinbox/gohttp/controller"

	"net/http"
)

type DemoContext struct {
	*api.BaseApiContext
}

func (d *DemoContext) BeforeAction() {
	d.BaseApiContext.BeforeAction()
}

type DemoController struct {
	api.ApiController
}

func (c *DemoController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) gcontroller.ActionContext {
	context := new(DemoContext)
	context.BaseApiContext = c.ApiController.NewActionContext(req, respWriter).(*api.BaseApiContext)

	return context
}

func (c *DemoController) demoLogic() demo.Logic {
	return factory.DefaultLogicFactory.DemoLogic()
}
