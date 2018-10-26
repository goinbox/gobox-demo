package mongodemo

import (
	"gdemo/controller/api"
	demoSvc "gdemo/svc/mongodemo"

	gcontroller "github.com/goinbox/gohttp/controller"

	"net/http"
)

type MongoDemoContext struct {
	*api.ApiContext

	demoSvc *demoSvc.MongoDemoSvc
}

func (d *MongoDemoContext) BeforeAction() {
	d.ApiContext.BeforeAction()

	d.demoSvc = demoSvc.NewMongoDemoSvc(
		d.AccessLogger,
		d.MongoClient,
	)
}

type MongoDemoController struct {
	api.BaseController
}

func (d *MongoDemoController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) gcontroller.ActionContext {
	context := new(MongoDemoContext)
	context.ApiContext = d.BaseController.NewActionContext(req, respWriter).(*api.ApiContext)

	return context
}
