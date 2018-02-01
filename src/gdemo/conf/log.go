package conf

import "time"

var LogConf struct {
	RootPath                    string
	Level                       int
	MaxAsyncMsgNum              int
	MaxBufferNum                int
	Bufsize                     int
	BufferAutoFlushTimeInterval time.Duration
}

func initLogConf() {
	LogConf.RootPath = PrjHome + "/logs"
	LogConf.Level = scJson.Log.Level
	LogConf.MaxAsyncMsgNum = scJson.Log.MaxAsyncMsgNum
	LogConf.MaxBufferNum = scJson.Log.MaxBufferNum
	LogConf.Bufsize = scJson.Log.Bufsize
	LogConf.BufferAutoFlushTimeInterval = time.Duration(scJson.Log.BufferAutoFlushTimeIntervalSeconds) * time.Second
}
