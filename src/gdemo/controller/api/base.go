package api

import (
	"gdemo/controller"
	"gdemo/errno"
	"gdemo/gvalue"
	"gdemo/misc"
	"github.com/goinbox/gohttp/system"
	"github.com/goinbox/mysql"

	"github.com/goinbox/exception"
	gcontroller "github.com/goinbox/gohttp/controller"
	"github.com/goinbox/golog"
	"github.com/goinbox/redis"

	"html"
	"net/http"
)

type ApiContext struct {
	*controller.BaseContext

	ApiData struct {
		V    string
		Data interface{}
		Err  *exception.Exception
	}

	MysqlClient *mysql.Client
	MysqlLogger golog.ILogger

	RedisPool   *redis.Pool
	RedisClient *redis.Client
	RedisLogger golog.ILogger
}

func (a *ApiContext) BeforeAction() {
	a.BaseContext.BeforeAction()

	var err error
	a.MysqlClient, err = gvalue.NewMysqlClient()
	if err != nil {
		a.ApiData.Err = exception.New(errno.E_SYS_MYSQL_ERROR, err.Error())
		system.JumpOutAction(JumpToApiError)
	}
	a.MysqlLogger = gvalue.NewAsyncLogger(gvalue.MysqlLogWriter, a.LogFormater)
	a.MysqlClient.SetLogger(a.MysqlLogger)

	a.RedisPool = gvalue.RedisClientPool
	a.RedisClient, err = a.RedisPool.Get()
	if err != nil {
		a.ApiData.Err = exception.New(errno.E_SYS_REDIS_ERROR, err.Error())
		system.JumpOutAction(JumpToApiError)
	}
	a.RedisLogger = gvalue.NewAsyncLogger(gvalue.RedisLogWriter, a.LogFormater)
	a.RedisClient.SetLogger(a.RedisLogger)
}

func (a *ApiContext) AfterAction() {
	f := a.QueryValues.Get("fmt")
	if f == "jsonp" {
		callback := a.QueryValues.Get("_callback")
		if callback != "" {
			a.RespBody = misc.ApiJsonp(a.ApiData.V, a.ApiData.Data, a.ApiData.Err, html.EscapeString(callback))
			return
		}
	}

	a.RespBody = misc.ApiJson(a.ApiData.V, a.ApiData.Data, a.ApiData.Err)

	a.BaseContext.AfterAction()
}

func (a *ApiContext) Destruct() {
	a.MysqlClient.Free()
	a.MysqlLogger.Free()

	if a.RedisClient.Connected() {
		a.RedisClient.SetLogger(gvalue.NoopLogger)
		a.RedisPool.Put(a.RedisClient)
	}
	a.RedisLogger.Free()

	a.BaseContext.Destruct()
}

type BaseController struct {
	controller.BaseController
}

func (b *BaseController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) gcontroller.ActionContext {
	context := new(ApiContext)
	context.BaseContext = b.BaseController.NewActionContext(req, respWriter).(*controller.BaseContext)

	return context
}

func JumpToApiError(context gcontroller.ActionContext, args ...interface{}) {
	acontext := context.(*ApiContext)

	acontext.RespBody = misc.ApiJson(acontext.ApiData.V, acontext.ApiData.Data, acontext.ApiData.Err)
}
