package errno

const (
	Success = 0

	ESysInvalidPrjHome     = 11
	ESysInitServerConfFail = 12
	ESysInitLogFail        = 13
	ESysSavePidFileFail    = 14
	ESysMysqlError         = 15
	ESysRedisError         = 16
	ESysMongoError         = 17

	ECommonFileNotExist            = 101
	ECommonReadFileError           = 102
	ECommonJsonEncodeError         = 103
	ECommonJsonDecodeError         = 104
	ECommonInvalidApiFmt           = 105
	ECommonInvalidApiJsonpCallback = 106
	ECommonInvalidArg              = 107
	ECommonInsertEntityFailed      = 108
	ECommonUpdateEntityFailed      = 109
)
