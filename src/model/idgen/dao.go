package idgen

import (
	"gdemo/pcontext"
	"github.com/goinbox/mysql"

	"gdemo/model"
)

const (
	genIDSql = "UPDATE id_gen SET max_id = last_insert_id(max_id + 1) WHERE name = ?"
)

type Dao interface {
	model.Dao

	GenID(ctx *pcontext.Context, name string) (int64, error)
}

type dao struct {
	*model.BaseDao
}

func NewDao(client *mysql.Client) Dao {
	return &dao{
		BaseDao: model.NewBaseDao("id_gen", client),
	}
}

func (d *dao) GenID(ctx *pcontext.Context, name string) (int64, error) {
	id, err := d.genID(ctx, name)
	if id != 0 {
		return id, nil
	}

	if err != nil {
		return 0, err
	}

	err = d.Insert(ctx, &Entity{
		Name:  name,
		MaxID: 0,
	}).Err

	if err != nil {
		if !mysql.DuplicateError(err) {
			return 0, err
		}
	}

	return d.genID(ctx, name)
}

func (d *dao) genID(ctx *pcontext.Context, name string) (int64, error) {
	result, err := d.Dao.Exec(ctx, genIDSql, name)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
