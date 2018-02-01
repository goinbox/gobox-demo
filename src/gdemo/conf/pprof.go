package conf

var PprofConf struct {
	Enable bool
	Port   string
}

func initPprofConf() {
	if scJson.Pprof.OnlineEnableHost == BaseConf.Hostname {
		PprofConf.Enable = true
	} else {
		PprofConf.Enable = false
	}

	PprofConf.Port = scJson.Pprof.Port
}
