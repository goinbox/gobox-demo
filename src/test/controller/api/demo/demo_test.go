package demo

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goinbox/gohttp/httpserver"
	"github.com/goinbox/gohttp/router"

	"gdemo/controller/api/demo"
	"gdemo/logic/factory"
	"gdemo/test"
)

var runner *test.ApiControllerRunner

func init() {
	dir, _ := os.Getwd()
	for i := 0; i < 5; i++ {
		dir = filepath.Dir(dir)
	}

	test.InitTestResource(dir)

	r := router.NewRouter()
	r.MapRouteItems(new(demo.Controller))

	runner = &test.ApiControllerRunner{
		Server: httpserver.NewServer(r),
		App:    factory.DefaultLogicFactory.AppLogic().ListAllApps(test.Context())[0],
	}
}

func TestAdd(t *testing.T) {
	content, err := runner.Run("/Demo/Add", `
{
  "Name": "b"
}
`)
	t.Log(err, string(content))
}

func TestIndex(t *testing.T) {
	content, err := runner.Run("/Demo/Index", `
{
  "Status": 1,
  "Offset": 0,
  "Limit": 10,
  "OrderBy": "name",
  "Order": "asc"
}
`)
	t.Log(err, string(content))

}

func TestEdit(t *testing.T) {
	content, err := runner.Run("/Demo/Edit", `
{
  "ID": 21,
  "Name": "b",
  "Status": 0
}
`)
	t.Log(err, string(content))
}

func TestDel(t *testing.T) {
	content, err := runner.Run("/Demo/Del", `
{
  "IDs": [
    13,27
  ]
}
`)
	t.Log(err, string(content))
}
