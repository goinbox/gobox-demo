package conf

type logConfJson struct {
	Level          int `json:"level"`
	AsyncQueueSize int `json:"async_queue_size"`
	Bufsize        int `json:"bufsize"`
}

var LogConf struct {
	RootPath       string
	Level          int
	AsyncQueueSize int
	Bufsize        int
}

func initLogConf() {
	LogConf.RootPath = PrjHome + "/logs"
	LogConf.Level = scJson.Log.Level
	LogConf.AsyncQueueSize = scJson.Log.AsyncQueueSize
	LogConf.Bufsize = scJson.Log.Bufsize
}
