package conf

import (
	"github.com/goinbox/exception"
	"github.com/goinbox/gomisc"

	"gdemo/errno"
)

var PrjHome string

func Init(prjHome string) *exception.Exception {
	if !gomisc.DirExist(prjHome) {
		return exception.New(errno.E_SYS_INVALID_PRJ_HOME, "prjHome not exists")
	}

	PrjHome = prjHome

	err := initServerConfJson()
	if err != nil {
		return exception.New(errno.E_SYS_INIT_SERVER_CONF_JSON_FAIL, "init serverConfJson error: "+err.Error())
	}

	initBaseConf()
	initLogConf()
	initPprofConf()
	initHttpConf()
	initRedisConf()
	initMysqlConf()

	return nil
}
