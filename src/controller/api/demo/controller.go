package demo

import (
	"gdemo/controller/api"
	"gdemo/logic/demo"
	"gdemo/logic/factory"

	gcontroller "github.com/goinbox/gohttp/controller"

	"net/http"
)

type DemoContext struct {
	*api.ApiContext
}

func (d *DemoContext) BeforeAction() {
	d.ApiContext.BeforeAction()
}

type DemoController struct {
	api.ApiController
}

func (c *DemoController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) gcontroller.ActionContext {
	context := new(DemoContext)
	context.ApiContext = c.ApiController.NewActionContext(req, respWriter).(*api.ApiContext)

	return context
}

func (c *DemoController) demoLogic() demo.Logic {
	return factory.DefaultLogicFactory.DemoLogic()
}
