package conf

import (
	"github.com/goinbox/goerror"
	"github.com/goinbox/gomisc"

	"gdemo/errno"
)

var PrjHome string

func Init(prjHome string) *goerror.Error {
	if !gomisc.DirExist(prjHome) {
		return goerror.New(errno.ESysInvalidPrjHome, "prjHome not exists")
	}

	PrjHome = prjHome

	err := initServerConfJson()
	if err != nil {
		return goerror.New(errno.ESysInitServerConfFail, "init serverConfJson error: "+err.Error())
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
