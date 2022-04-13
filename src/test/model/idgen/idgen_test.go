package idgen

import (
	"os"
	"path/filepath"
	"testing"

	"gdemo/model/idgen"
	"gdemo/test"
)

func init() {
	dir, _ := os.Getwd()
	for i := 0; i < 4; i++ {
		dir = filepath.Dir(dir)
	}

	test.InitTestResource(dir)
}

func TestGenID(t *testing.T) {
	dao := idgen.NewDao(test.MysqlClient())
	for i := 0; i < 10; i++ {
		id, err := dao.GenID("demo")
		t.Log(id, err)
	}
}
