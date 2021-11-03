package demo

import (
	"github.com/goinbox/mysql"

	"gdemo/model"
)

type Dao interface {
	model.Dao
}

type dao struct {
	*model.BaseDao
}

func NewDao(client *mysql.Client) Dao {
	return &dao{
		BaseDao: model.NewBaseDao("demo", client),
	}
}
