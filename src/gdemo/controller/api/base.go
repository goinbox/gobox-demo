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
	"github.com/goinbox/redis"
	"github.com/goinbox/golog"

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

func (this *ApiContext) BeforeAction() {
	this.BaseContext.BeforeAction()

	var err error
	this.MysqlClient, err = gvalue.NewMysqlClient()
	if err != nil {
		this.ApiData.Err = exception.New(errno.E_SYS_MYSQL_ERROR, err.Error())
		system.JumpOutAction(JumpToApiError)
	}
	this.MysqlLogger = gvalue.NewAsyncLogger(gvalue.MysqlLogWriter, this.LogFormater)
	this.MysqlClient.SetLogger(this.MysqlLogger)

	this.RedisPool = gvalue.RedisClientPool
	this.RedisClient, err = this.RedisPool.Get()
	if err != nil {
		this.ApiData.Err = exception.New(errno.E_SYS_REDIS_ERROR, err.Error())
		system.JumpOutAction(JumpToApiError)
	}
	this.RedisLogger = gvalue.NewAsyncLogger(gvalue.RedisLogWriter, this.LogFormater)
	this.RedisClient.SetLogger(this.RedisLogger)
}

func (this *ApiContext) AfterAction() {
	f := this.QueryValues.Get("fmt")
	if f == "jsonp" {
		callback := this.QueryValues.Get("_callback")
		if callback != "" {
			this.RespBody = misc.ApiJsonp(this.ApiData.V, this.ApiData.Data, this.ApiData.Err, html.EscapeString(callback))
			return
		}
	}

	this.RespBody = misc.ApiJson(this.ApiData.V, this.ApiData.Data, this.ApiData.Err)

	this.BaseContext.AfterAction()
}

func (this *ApiContext) Destruct() {
	this.MysqlClient.Free()
	this.MysqlLogger.Free()

	if this.RedisClient.Connected() {
		this.RedisClient.SetLogger(gvalue.NoopLogger)
		this.RedisPool.Put(this.RedisClient)
	}
	this.RedisLogger.Free()

	this.BaseContext.Destruct()
}

type BaseController struct {
	controller.BaseController
}

func (this *BaseController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) gcontroller.ActionContext {
	context := new(ApiContext)
	context.BaseContext = this.BaseController.NewActionContext(req, respWriter).(*controller.BaseContext)

	return context
}

func JumpToApiError(context gcontroller.ActionContext, args ...interface{}) {
	acontext := context.(*ApiContext)

	acontext.RespBody = misc.ApiJson(acontext.ApiData.V, acontext.ApiData.Data, acontext.ApiData.Err)
}
