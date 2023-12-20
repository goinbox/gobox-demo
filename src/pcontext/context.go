package pcontext

import (
	"context"
	"github.com/goinbox/gohttp/v6/httpserver"

	"github.com/goinbox/golog"
	"github.com/goinbox/mysql"
	"github.com/goinbox/pcontext"
)

type Context struct {
	httpserver.Context

	tid string

	mysqlClient *mysql.Client
}

func NewContext(logger golog.Logger) *Context {
	return &Context{
		Context: &httpserver.BaseContext{
			Context: pcontext.NewSimpleContext(nil, logger),
		},
	}
}

func (c *Context) TraceID() string {
	return c.tid
}

func (c *Context) SetTraceID(tid string) *Context {
	c.tid = tid

	return c
}

func (c *Context) MySQLClient() *mysql.Client {
	return c.mysqlClient
}

func (c *Context) SetMySQLClient(client *mysql.Client) *Context {
	c.mysqlClient = client

	return c
}

func (c *Context) WithContext(ctx context.Context) *Context {
	return c.copy(ctx)
}

func (c *Context) copy(ctx context.Context) *Context {
	cc := &Context{
		Context: &httpserver.BaseContext{
			Context: pcontext.NewSimpleContext(ctx, c.Logger()),
		},
	}

	cc.SetController(c.Controller())
	cc.SetAction(c.Action())

	cc.tid = c.tid
	cc.mysqlClient = c.mysqlClient

	return cc
}
