package api

import (
	"gdemo/conf"
	"gdemo/controller"
	"gdemo/errno"
	"gdemo/gvalue"
	"gdemo/misc"

	"github.com/goinbox/exception"
	gcontroller "github.com/goinbox/gohttp/controller"
	"github.com/goinbox/gohttp/query"
	"github.com/goinbox/gohttp/system"
	"github.com/goinbox/golog"
	"github.com/goinbox/mongo"
	"github.com/goinbox/mysql"
	"github.com/goinbox/redis"

	"html"
	"net/http"
	"net/url"
	"time"
)

type ApiSignParams struct {
	T     int64
	Nonce string
	Sign  string
	Debug int
}

var ApiSignQueryNames = []string{"t", "nonce"}

func SetApiSignParams(qs *query.QuerySet, asp *ApiSignParams) {
	qs.Int64Var(&asp.T, "t", true, errno.E_COMMON_INVALID_ARG, "invalid sign t", query.CheckInt64IsPositive)
	qs.StringVar(&asp.Nonce, "nonce", true, errno.E_COMMON_INVALID_ARG, "invalid sign nonce", query.CheckStringNotEmpty)
	qs.StringVar(&asp.Sign, "sign", true, errno.E_COMMON_INVALID_ARG, "invalid sign sign", query.CheckStringNotEmpty)
	qs.IntVar(&asp.Debug, "debug", false, errno.E_COMMON_INVALID_ARG, "invalid sign debug", nil)
}

func VerifyApiSign(asp *ApiSignParams, queryValues url.Values, signQueryNames []string, token string) *exception.Exception {
	if conf.BaseConf.IsDev && asp.Debug == 1 {
		return nil
	}

	if time.Now().Unix()-asp.T > 600 {
		return exception.New(errno.E_COMMON_INVALID_ARG, "verify sign failed, invalid sign t")
	}

	sign := misc.CalApiSign(queryValues, signQueryNames, token)
	if sign != asp.Sign {
		return exception.New(errno.E_COMMON_INVALID_ARG, "verify sign failed, invalid sign sign")
	}

	return nil
}

type IApiDataContext interface {
	gcontroller.ActionContext

	Version() string
	Data() interface{}
	Err() *exception.Exception
}

type ApiContext struct {
	*controller.BaseContext

	ApiData struct {
		V    string
		Data interface{}
		Err  *exception.Exception
	}

	MysqlClient *mysql.Client

	RedisPool   *redis.Pool
	RedisClient *redis.Client

	MongoPool   *mongo.Pool
	MongoClient *mongo.Client
	MongoLogger golog.ILogger
}

func (a *ApiContext) Version() string {
	return a.ApiData.V
}

func (a *ApiContext) Data() interface{} {
	return a.ApiData.Data
}

func (a *ApiContext) Err() *exception.Exception {
	return a.ApiData.Err
}

func (a *ApiContext) BeforeAction() {
	a.BaseContext.BeforeAction()

	var err error
	a.MysqlClient, err = gvalue.NewMysqlClient()
	if err != nil {
		a.ApiData.Err = exception.New(errno.E_SYS_MYSQL_ERROR, err.Error())
		system.JumpOutAction(JumpToApiError)
	}
	a.MysqlClient.SetLogger(a.AccessLogger)

	a.RedisPool = gvalue.RedisClientPool
	a.RedisClient, err = a.RedisPool.Get()
	if err != nil {
		a.ApiData.Err = exception.New(errno.E_SYS_REDIS_ERROR, err.Error())
		system.JumpOutAction(JumpToApiError)
	}
	a.RedisClient.SetLogger(a.AccessLogger)

	a.MongoPool = gvalue.MongoClientPool
	a.MongoClient, err = a.MongoPool.Get()
	if err != nil {
		a.ApiData.Err = exception.New(errno.E_SYS_MONGO_ERROR, err.Error())
		system.JumpOutAction(JumpToApiError)
	}
	a.MongoLogger = gvalue.NewAsyncLogger(gvalue.MongoLogWriter, a.LogFormater)
	a.MongoClient.SetLogger(a.MongoLogger)
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
	if a.MysqlClient != nil {
		a.MysqlClient.Free()
	}

	if a.RedisClient != nil {
		if a.RedisClient.Connected() {
			a.RedisClient.SetLogger(gvalue.NoopLogger)
			a.RedisPool.Put(a.RedisClient)
		}
	}

	if a.MongoClient != nil {
		if a.MongoClient.Connected() {
			a.MongoClient.SetLogger(gvalue.NoopLogger)
			a.MongoPool.Put(a.MongoClient)
		}
		if a.MongoLogger != nil {
			a.MongoLogger.Free()
		}
	}

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
	acontext := context.(IApiDataContext)

	acontext.SetResponseBody(misc.ApiJson(acontext.Version(), acontext.Data(), acontext.Err()))
}
