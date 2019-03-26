package conf

import (
	"github.com/goinbox/gomisc"
)

var scJson serverConfJson

type serverConfJson struct {
	PrjName string `json:"prj_name"`
	IsDev   bool   `json:"is_dev"`
	Idc     string `json:"idc"`

	Log     logConfJson   `json:"log"`
	Pprof   pprofConfJson `json:"pprof"`
	ApiHttp httpConfJson  `json:"api_http"`

	Redis []*redisConfJson `json:"redis"`
	Mysql mysqlConfJson    `json:"mysql"`
	Mongo mongoConfJson    `json:"mongo"`
}

func initServerConfJson() error {
	confRoot := PrjHome + "/conf"
	err := gomisc.ParseJsonFile(confRoot+"/server/server_conf.json", &scJson)
	if err != nil {
		return err
	}
	err = gomisc.ParseJsonFile(confRoot+"/server_conf_rewrite.json", &scJson)
	if err != nil {
		return err
	}

	return nil
}
