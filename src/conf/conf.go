package conf

import (
	"github.com/goinbox/exception"
	"github.com/goinbox/gomisc"

	"gdemo/errno"
)

var PrjHome string

func Init(prjHome string) *exception.Exception {
	if !gomisc.DirExist(prjHome) {
		return exception.New(errno.ESysInvalidPrjHome, "prjHome not exists")
	}

	PrjHome = prjHome

	err := initServerConfJson()
	if err != nil {
		return exception.New(errno.ESysInitServerConfFail, "init serverConfJson error: "+err.Error())
	}

	initBaseConf()
	initLogConf()
	initPprofConf()
	initHttpConf()
	initRedisConf()
	initMysqlConf()
	initMongoConf()

	return nil
}
