package demo

import (
	"testing"

	"gdemo/model"
	"gdemo/model/demo"
	"gdemo/test"
)

func init() {
	test.InitMysql()
}

func TestDemoDao(t *testing.T) {
	entity := &demo.Entity{
		BaseEntity: model.BaseEntity{},
		Name:       "",
		Status:     0,
	}
}

func dao() demo.Dao {
	return demo.NewDao(test.MysqlClient())
}
