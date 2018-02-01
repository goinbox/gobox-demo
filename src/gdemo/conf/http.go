package conf

var ApiHttpConf struct {
	GoHttpHost string
	GoHttpPort string
}

func initHttpConf() {
	ApiHttpConf.GoHttpHost = scJson.ApiHttp.GoHttpHost
	ApiHttpConf.GoHttpPort = scJson.ApiHttp.GoHttpPort
}
