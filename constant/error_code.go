package constant

const (
	PARAM_ERROR                 uint64 = 1001 // 参数错误
	JWT_VERIFY_ERROR            uint64 = 1002 // jwt权限验证失败
	API_RATELIMIT               uint64 = 1003 // api频率限制
	AUTH_ERROR                  uint64 = 1004 // api权限验证不通过
	ILLEGAL_ADDRESS             uint64 = 1005 // 地址非法
	USER_CURRENCY_NOT_AVAILABLE uint64 = 1006 // 用户币种不可用
	CURRENCY_NOT_AVAILABLE      uint64 = 1009 // 币种不可用（无此币种或币种被禁用）

	WITHDRAW_REQUEST_ID_USED uint64 = 1007 // 提现request id已经被使用
	CURRENCY_WITHDRAW_BANNED uint64 = 1008 // 币种提现被禁用
	AMOUNT_DECIMAL_ERR       uint64 = 1009 // 数量精度有问题
	MIN_WITHDRAW_AMOUNT      uint64 = 1010 // 未达到币种最小提现金额
	USER_MAX_WITHDRAW_AMOUNT uint64 = 1011 // 超过了最大提币金额
	BALANCE_NOT_ENOUGH       uint64 = 1012 // 余额不足
	WITHDRAW_LIMIT_DAILY     uint64 = 1013 // 超过单日限额
	MEMO_TOO_LONG            uint64 = 1014 // memo太长
)
