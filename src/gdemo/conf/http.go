package conf

type httpConfJson struct {
	GoHttpHost string `json:"gohttp_host"`
	GoHttpPort string `json:"gohttp_port"`
}

var ApiHttpConf struct {
	GoHttpHost string
	GoHttpPort string
}

func initHttpConf() {
	ApiHttpConf.GoHttpHost = scJson.ApiHttp.GoHttpHost
	ApiHttpConf.GoHttpPort = scJson.ApiHttp.GoHttpPort
}
