package conf

import (
	"os"
	"os/user"
)

var BaseConf struct {
	Hostname string
	Username string

	PrjName string
	IsDev   bool
	Idc     string

	TmpRoot    string
	ApiPidFile string
}

func initBaseConf() {
	BaseConf.Hostname, _ = os.Hostname()
	curUser, _ := user.Current()
	BaseConf.Username = curUser.Username

	BaseConf.PrjName = scJson.PrjName
	BaseConf.IsDev = scJson.IsDev
	BaseConf.Idc = scJson.Idc

	BaseConf.TmpRoot = PrjHome + "/tmp"
	BaseConf.ApiPidFile = BaseConf.TmpRoot + "/api.pid"
}
