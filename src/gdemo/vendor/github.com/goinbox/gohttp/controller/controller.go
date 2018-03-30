package controller

import (
	"net/http"
)

type ActionContext interface {
	Request() *http.Request
	ResponseWriter() http.ResponseWriter

	ResponseBody() []byte
	SetResponseBody(body []byte)

	BeforeAction()
	AfterAction()
	Destruct()
}

type Controller interface {
	NewActionContext(req *http.Request, respWriter http.ResponseWriter) ActionContext
}
