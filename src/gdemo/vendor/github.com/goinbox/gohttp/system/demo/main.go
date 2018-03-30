package main

import (
	"github.com/goinbox/gohttp/controller"
	"github.com/goinbox/gohttp/gracehttp"
	"github.com/goinbox/gohttp/router"
	"github.com/goinbox/gohttp/system"

	"net/http"
)

func main() {
	dcl := new(DemoController)
	r := router.NewSimpleRouter()

	r.DefineRouteItem("^/g/([0-9]+)$", dcl, "get")
	r.MapRouteItems(new(IndexController), dcl)

	sys := system.NewSystem(r)

	gracehttp.ListenAndServe(":8001", sys)
}

type BaseActionContext struct {
	Req        *http.Request
	RespWriter http.ResponseWriter
	RespBody   []byte
}

func (this *BaseActionContext) Request() *http.Request {
	return this.Req
}

func (this *BaseActionContext) ResponseWriter() http.ResponseWriter {
	return this.RespWriter
}

func (this *BaseActionContext) ResponseBody() []byte {
	return this.RespBody
}

func (this *BaseActionContext) SetResponseBody(body []byte) {
	this.RespBody = body
}

func (this *BaseActionContext) BeforeAction() {
	this.RespBody = append(this.RespBody, []byte(" index before ")...)
}

func (this *BaseActionContext) AfterAction() {
	this.RespBody = append(this.RespBody, []byte(" index after ")...)
}

func (this *BaseActionContext) Destruct() {
	println(" index destruct ")
}

type IndexController struct {
}

func (this *IndexController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) controller.ActionContext {
	return &BaseActionContext{
		Req:        req,
		RespWriter: respWriter,
	}
}

func (this *IndexController) IndexAction(context *BaseActionContext) {
	context.RespBody = append(context.RespBody, []byte(" index action ")...)
}

func (this *IndexController) RedirectAction(context *BaseActionContext) {
	system.Redirect302("https://github.com/goinbox")
}

type DemoActionContext struct {
	*BaseActionContext
}

func (this *DemoActionContext) BeforeAction() {
	this.RespBody = append(this.RespBody, []byte(" demo before ")...)
}

func (this *DemoActionContext) AfterAction() {
	this.RespBody = append(this.RespBody, []byte(" demo after ")...)
}

func (this *DemoActionContext) Destruct() {
	println(" demo destruct ")
}

type DemoController struct {
}

func (this *DemoController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) controller.ActionContext {
	return &DemoActionContext{
		&BaseActionContext{
			Req:        req,
			RespWriter: respWriter,
		},
	}
}

func (this *DemoController) DemoAction(context *DemoActionContext) {
	context.RespBody = append(context.RespBody, []byte(" demo action ")...)
}

func (this *DemoController) GetAction(context *DemoActionContext, id string) {
	context.RespBody = append(context.RespBody, []byte(" get action id = "+id)...)
}
