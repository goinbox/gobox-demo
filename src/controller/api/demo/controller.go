package demo

import (
	"net/http"
)

type Controller struct {
}

func (c *Controller) Name() string {
	return "demo"
}

func (c *Controller) IndexAction(r *http.Request, w http.ResponseWriter, args []string) *indexAction {
	return newIndexAction(r, w, args)
}

func (c *Controller) AddAction(r *http.Request, w http.ResponseWriter, args []string) *addAction {
	return newAddAction(r, w, args)
}

func (c *Controller) EditAction(r *http.Request, w http.ResponseWriter, args []string) *editAction {
	return newEditAction(r, w, args)
}

func (c *Controller) DelAction(r *http.Request, w http.ResponseWriter, args []string) *delAction {
	return newDelAction(r, w, args)
}
