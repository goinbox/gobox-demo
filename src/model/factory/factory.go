package factory

import (
	"github.com/goinbox/mysql"

	"gdemo/model/demo"
	"gdemo/model/idgen"
	"gdemo/pcontext"
	"gdemo/resource"
)

type DaoFactory interface {
	IDGenDao(ctx *pcontext.Context) idgen.Dao
	DemoDao(ctx *pcontext.Context) demo.Dao
}

var DefaultDaoFactory DaoFactory = new(daoFactory)

type daoFactory struct {
}

func (f *daoFactory) client(ctx *pcontext.Context) *mysql.Client {
	if ctx.MySQLClient == nil {
		ctx.MySQLClient = resource.MySQLClient(ctx.Logger)
	}

	return ctx.MySQLClient
}

func (f *daoFactory) IDGenDao(ctx *pcontext.Context) idgen.Dao {
	return idgen.NewDao(f.client(ctx))
}

func (f *daoFactory) DemoDao(ctx *pcontext.Context) demo.Dao {
	return demo.NewDao(f.client(ctx))
}
