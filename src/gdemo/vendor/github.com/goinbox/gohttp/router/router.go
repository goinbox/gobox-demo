package router

import (
	"reflect"

	"github.com/goinbox/gohttp/controller"
)

type Route struct {
	Cl          controller.Controller
	ActionValue *reflect.Value
	Args        []string
}

type Router interface {
	MapRouteItems(cls ...controller.Controller)
	DefineRouteItem(pattern string, cl controller.Controller, actionName string)

	FindRoute(path string) *Route
}
