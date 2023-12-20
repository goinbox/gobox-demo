package perror

const (
	Success = 0

	ESysInvalidConfDir  = 11
	ESysInitConfError   = 12
	ESysInitLogFail     = 13
	ESysFileIOError     = 14
	ESysMysqlError      = 15
	ESysRedisError      = 16
	ESysRunServerError  = 17
	ESysInitTracingFail = 18

	ECommonSysError         = 101
	ECommonInvalidArg       = 102
	ECommonDataAlreadyExist = 103
	ECommonAuthFailure      = 104
)
