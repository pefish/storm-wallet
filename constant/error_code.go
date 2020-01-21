package constant

const (
	PARAM_ERROR                 uint64 = 2001 // 参数错误
	API_RATELIMIT               uint64 = 2003 // api频率限制
	AUTH_ERROR                  uint64 = 2004 // api权限验证不通过
	ILLEGAL_ADDRESS             uint64 = 2005 // 地址非法
	USER_CURRENCY_NOT_AVAILABLE uint64 = 2006 // 用户币种不可用
	CURRENCY_NOT_AVAILABLE      uint64 = 2009 // 币种不可用（无此币种或币种被禁用）

	WITHDRAW_REQUEST_ID_USED uint64 = 2007 // 提现request id已经被使用
	CURRENCY_WITHDRAW_BANNED uint64 = 2008 // 币种提现被禁用
	AMOUNT_DECIMAL_ERR       uint64 = 2009 // 数量精度有问题
	MIN_WITHDRAW_AMOUNT      uint64 = 2010 // 未达到币种最小提现金额
	USER_MAX_WITHDRAW_AMOUNT uint64 = 2011 // 超过了最大提币金额
	BALANCE_NOT_ENOUGH       uint64 = 2012 // 余额不足
	WITHDRAW_LIMIT_DAILY     uint64 = 2013 // 超过单日限额
	MEMO_TOO_LONG            uint64 = 2014 // memo太长
	BALANCE_EXCEPTION        uint64 = 2015 // 资产异常
	ILLEGAL_USER             uint64 = 2016 // 用户非法
	WITHDRAW_CONFIRM_FAIL    uint64 = 2017 // 提现二次确认失败
	TX_NOT_FOUND             uint64 = 2018 // 交易没找到
	CURRENCY_DEPOSIT_BANNED  uint64 = 2008 // 币种充值被禁用
	ADDRESS_INDEX_TOO_BIG    uint64 = 2009 // 地址索引太大
	USERID_TOO_BIG           uint64 = 2010 // 用户id太大
	WALLET_CONFIG_ERROR      uint64 = 2011 // 钱包配置错误
	ILLEGAL_MEMBER           uint64 = 2012 // 成员非法

	ROLE_AUTH_ERROR     uint64 = 2101 // 用户角色权限不够
	USER_NOT_FOUND      uint64 = 2102 // 用户没找到或被禁用
	NO_TEAM_ERROR       uint64 = 2103 // 没有加入团队
	ING_TEAM            uint64 = 2104 // 已经处于团队当中
	USER_NOT_IN_MY_TEAM uint64 = 2105 // 此成员不是本团队成员
)
