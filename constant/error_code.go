package constant

const (
	PARAM_ERROR      uint64 = 1001 // 参数错误
	JWT_VERIFY_ERROR uint64 = 1002 // jwt权限验证失败
	API_RATELIMIT    uint64 = 1003 // api频率限制
	AUTH_ERROR       uint64 = 1004 // api权限验证不通过
	ILLEGAL_ADDRESS  uint64 = 1005 // 地址非法
)
