package demo

import (
	"os"
	"path/filepath"
	"testing"

	"gdemo/controller/api/demo"
	"gdemo/test"
	"gdemo/test/controller/api"
)

func init() {
	dir, _ := os.Getwd()
	for i := 0; i < 5; i++ {
		dir = filepath.Dir(dir)
	}

	test.InitTestResource(dir)

	api.InitTestServer(new(demo.Controller))
}

func TestAdd(t *testing.T) {
	api.HandleRequest("/Demo/Add", `
{
  "Name": "b"
}
`)
}

func TestIndex(t *testing.T) {
	api.HandleRequest("/Demo/Index", `
{
  "Status": 1
}
`)
}

func TestEdit(t *testing.T) {
	api.HandleRequest("/Demo/Edit", `
{
  "ID": 13,
  "Status": 0
}
`)
}

func TestDel(t *testing.T) {
	api.HandleRequest("/Demo/Del", `
{
  "IDs": [
    13
  ]
}
`)
}
