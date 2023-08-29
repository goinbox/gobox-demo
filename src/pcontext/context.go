package pcontext

import (
	"github.com/goinbox/golog"
	"github.com/goinbox/mysql"
	"github.com/goinbox/pcontext"
)

type Context struct {
	pcontext.Context

	tid string

	mysqlClient *mysql.Client
}

func NewContext(logger golog.Logger, tid string) *Context {
	return &Context{
		Context: pcontext.NewSimpleContext(logger),
		tid:     tid,
	}
}

func (c *Context) TraceID() string {
	return c.tid
}

func (c *Context) MySQLClient() *mysql.Client {
	return c.mysqlClient
}

func (c *Context) SetMySQLClient(client *mysql.Client) *Context {
	c.mysqlClient = client

	return c
}
