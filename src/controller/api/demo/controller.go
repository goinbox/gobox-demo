package demo

import (
	"net/http"

	"github.com/goinbox/gohttp/httpserver"
)

type Controller struct {
}

func (c *Controller) Name() string {
	return "Demo"
}

func (c *Controller) IndexAction(r *http.Request, w http.ResponseWriter, args []string) httpserver.Action {
	return newIndexAction(r, w, args)
}

func (c *Controller) AddAction(r *http.Request, w http.ResponseWriter, args []string) httpserver.Action {
	return newAddAction(r, w, args)
}

func (c *Controller) EditAction(r *http.Request, w http.ResponseWriter, args []string) httpserver.Action {
	return newEditAction(r, w, args)
}

func (c *Controller) DelAction(r *http.Request, w http.ResponseWriter, args []string) httpserver.Action {
	return newDelAction(r, w, args)
}
