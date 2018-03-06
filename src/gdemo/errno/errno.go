package errno

const (
	SUCCESS = 0

	E_SYS_INVALID_PRJ_HOME           = 11
	E_SYS_INIT_SERVER_CONF_JSON_FAIL = 12
	E_SYS_INIT_SERVER_CONF_FAIL      = 13
	E_SYS_INIT_LOG_FAIL              = 14
	E_SYS_SAVE_PID_FILE_FAIL         = 15
	E_SYS_MYSQL_ERROR                = 16
	E_SYS_REDIS_ERROR                = 17

	E_COMMON_FILE_NOT_EXIST             = 101
	E_COMMON_READ_FILE_ERROR            = 102
	E_COMMON_JSON_ENCODE_ERROR          = 103
	E_COMMON_JSON_DECODE_ERROR          = 104
	E_COMMON_INVALID_API_FMT            = 105
	E_COMMON_INVALID_API_JSONP_CALLBACK = 106
	E_COMMON_INVALID_ID                 = 107
	E_COMMON_INVALID_ADD_TIME           = 108
	E_COMMON_INVALID_EDIT_TIME          = 109
	E_COMMON_INVALID_QUERY_OFFSET       = 110
	E_COMMON_INVALID_QUERY_CNT          = 111
	E_COMMON_INVALID_SIGN_T             = 112
	E_COMMON_INVALID_SIGN_NONCE         = 113
	E_COMMON_INVALID_SIGN_SIGN          = 114
	E_COMMON_INVALID_SIGN_DEBUG         = 115

	E_API_DEMO_INVALID_NAME   = 1001
	E_API_DEMO_INVALID_STATUS = 1002
	E_API_DEMO_INSERT_FAILED  = 1003
	E_API_DEMO_UPDATE_FAILED  = 1004
)
