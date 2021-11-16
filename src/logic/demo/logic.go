package demo

import (
	"github.com/goinbox/mysql"

	"gdemo/model/demo"
	"gdemo/model/factory"
	"gdemo/pcontext"
)

type Logic interface {
	Insert(ctx *pcontext.Context, entities ...*demo.Entity) error
	DeleteByIDs(ctx *pcontext.Context, ids ...int64) error
	UpdateByIDs(ctx *pcontext.Context, fields map[string]interface{}, ids ...int64) error
	SelectByID(ctx *pcontext.Context, id int64) (*demo.Entity, error)

	SimpleQueryAnd(ctx *pcontext.Context, params *mysql.SqlQueryParams) ([]*demo.Entity, error)
	SimpleTotalAnd(ctx *pcontext.Context, items ...*mysql.SqlColQueryItem) (int64, error)
}

type logic struct {
}

func NewLogic() *logic {
	return &logic{}
}

func (l *logic) dao(ctx *pcontext.Context) demo.Dao {
	return factory.DefaultDaoFactory.DemoDao(ctx)
}

func (l *logic) Insert(ctx *pcontext.Context, entities ...*demo.Entity) error {
	data := make([]interface{}, len(entities))
	for i, entity := range entities {
		data[i] = entity
	}

	return l.dao(ctx).Insert(data...).Err
}

func (l *logic) DeleteByIDs(ctx *pcontext.Context, ids ...int64) error {
	return l.dao(ctx).DeleteByIDs(ids...).Err
}

func (l *logic) UpdateByIDs(ctx *pcontext.Context, fields map[string]interface{}, ids ...int64) error {
	return l.dao(ctx).UpdateByIDs(fields, ids...).Err
}

func (l *logic) SelectByID(ctx *pcontext.Context, id int64) (*demo.Entity, error) {
	entity := new(demo.Entity)

	err := l.dao(ctx).SelectByID(id, entity)

	return entity, err
}

func (l *logic) SimpleQueryAnd(ctx *pcontext.Context, params *mysql.SqlQueryParams) ([]*demo.Entity, error) {
	var entities []*demo.Entity

	err := l.dao(ctx).SimpleQueryAnd(params, &entities)

	return entities, err
}

func (l *logic) SimpleTotalAnd(ctx *pcontext.Context, items ...*mysql.SqlColQueryItem) (int64, error) {
	return l.dao(ctx).SimpleTotalAnd(items...)
}
