package conf

import (
	"gdemo/misc"
)

var scJson serverConfJson

type logConfJson struct {
	Level                              int `json:"level"`
	MaxAsyncMsgNum                     int `json:"max_async_msg_num"`
	MaxBufferNum                       int `json:"max_buffer_num"`
	Bufsize                            int `json:"bufsize"`
	BufferAutoFlushTimeIntervalSeconds int `json:"buffer_auto_flush_time_interval_seconds"`
}

type pprofConfJson struct {
	OnlineEnableHost string `json:"online_enable_host"`
	Port             string `json:"port"`
}

type httpConfJson struct {
	GoHttpHost string `json:"gohttp_host"`
	GoHttpPort string `json:"gohttp_port"`
}

type redisConfJson struct {
	Host                         string `json:"host"`
	Pass                         string `json:"pass"`
	Port                         string `json:"port"`
	RWTimeoutSeconds             int    `json:"rw_timeout_seconds"`
	PoolSize                     int    `json:"pool_size"`
	PoolKeepAliveIntervalSeconds int    `json:"pool_keepalive_interval_seconds"`
	PoolClientMaxIdleSeconds     int    `json:"pool_client_max_idle_seconds"`
}

type mysqlConfJson struct {
	Host                     string `json:"host"`
	User                     string `json:"user"`
	Pass                     string `json:"pass"`
	Port                     string `json:"port"`
	Name                     string `json:"name"`
	RWTimeoutSeconds         int    `json:"rw_timeout_seconds"`
	PoolSize                 int    `json:"pool_size"`
	PoolClientMaxIdleSeconds int    `json:"pool_client_max_idle_seconds"`
}

type serverConfJson struct {
	PrjName string `json:"prj_name"`
	IsDev   bool   `json:"is_dev"`
	Idc     string `json:"idc"`

	Log     logConfJson   `json:"log"`
	Pprof   pprofConfJson `json:"pprof"`
	ApiHttp httpConfJson  `json:"api_http"`

	Redis redisConfJson `json:"redis"`
	Mysql mysqlConfJson `json:"mysql"`
}

func initServerConfJson() error {
	confRoot := PrjHome + "/conf"
	err := misc.ParseJsonFile(confRoot+"/server/server_conf.json", &scJson)
	if err != nil {
		return err
	}
	err = misc.ParseJsonFile(confRoot+"/server_conf_rewrite.json", &scJson)
	if err != nil {
		return err
	}

	return nil
}
