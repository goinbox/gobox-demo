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
	client := ctx.MySQLClient()
	if client == nil {
		client = resource.MySQLClient()
		ctx.SetMySQLClient(client)
	}

	return client
}

func (f *daoFactory) IDGenDao(ctx *pcontext.Context) idgen.Dao {
	return idgen.NewDao(f.client(ctx))
}

func (f *daoFactory) DemoDao(ctx *pcontext.Context) demo.Dao {
	return demo.NewDao(f.client(ctx))
}
