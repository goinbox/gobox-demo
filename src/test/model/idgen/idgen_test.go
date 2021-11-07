package idgen

import (
	"testing"

	"gdemo/model/idgen"
	"gdemo/test"
)

func init() {
	test.InitMysql()
}

func TestGenID(t *testing.T) {
	dao := idgen.NewDao(test.MysqlClient())
	for i := 0; i < 10; i++ {
		id, err := dao.GenID("demo")
		t.Log(id, err)
	}
}
