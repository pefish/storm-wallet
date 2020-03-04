package controller

import (
	"fmt"
	"time"
	"wallet-storm-wallet/constant"
	external_service "wallet-storm-wallet/external-service"
	"wallet-storm-wallet/model"

	api_session "github.com/pefish/go-core/api-session"
	go_decimal "github.com/pefish/go-decimal"
	go_error "github.com/pefish/go-error"
	go_http "github.com/pefish/go-http"
	go_json "github.com/pefish/go-json"
	go_logger "github.com/pefish/go-logger"
	go_mysql "github.com/pefish/go-mysql"
	go_redis "github.com/pefish/go-redis"
	go_reflect "github.com/pefish/go-reflect"
	"github.com/pefish/storm-golang-sdk/signature"
	uuid "github.com/satori/go.uuid"
)

type WithdrawControllerClass struct {
}

var WithdrawController = WithdrawControllerClass{}

type WithdrawParam struct {
	Currency  string  `json:"currency" validate:"required" desc:"提现的币种名"`
	Chain     string  `json:"chain" validate:"required" desc:"提现的链名"`
	RequestId string  `json:"request_id" validate:"required,max=200" desc:"订单id。此id幂等"`
	Address   string  `json:"address" validate:"required" desc:"提现的目标地址"`
	Amount    string  `json:"amount" validate:"required" desc:"提现的数量"`
	Memo      *string `json:"memo,omitempty" validate:"omitempty" desc:"memo"`
}

func (this *WithdrawControllerClass) Withdraw(apiSession *api_session.ApiSessionClass) interface{} {
	params := WithdrawParam{}
	apiSession.ScanParams(&params)

	lockKey := fmt.Sprintf(`storm_wallet_withdraw_%d_lock`, apiSession.UserId)
	uniqueId := uuid.NewV1().String()
	if !go_redis.RedisHelper.MustGetLock(lockKey, uniqueId, 3*time.Second) {
		go_error.Throw(`rate limit`, constant.API_RATELIMIT)
	}
	defer go_redis.RedisHelper.MustReleaseLock(lockKey, uniqueId)

	// 检查request id是否已存在
	withdrawModel := model.WithdrawModel.GetByUserIdRequestId(apiSession.UserId, params.RequestId)
	if withdrawModel != nil {
		go_error.Throw(`request id is used already`, constant.WITHDRAW_REQUEST_ID_USED)
	}

	// 检查币种是否开启提币
	currencyModel := model.CurrencyModel.GetByCurrencyChain(params.Currency, params.Chain)
	if currencyModel == nil {
		go_error.Throw(`currency is not available`, constant.CURRENCY_NOT_AVAILABLE)
	}
	if currencyModel.IsWithdrawEnable == 0 {
		go_error.Throw(`currency withdraw banned`, constant.CURRENCY_WITHDRAW_BANNED)
	}

	// 检查用户是否具有此币种
	userCurrencyModel := model.UserCurrencyModel.GetByUserIdCurrencyId(apiSession.UserId, currencyModel.Id)
	if userCurrencyModel == nil {
		go_error.Throw(`user currency is not available`, constant.USER_CURRENCY_NOT_AVAILABLE)
	}

	// 校验数额与币种精度的匹配
	if go_decimal.Decimal.Start(params.Amount).GetPrecision() > int32(currencyModel.Decimals) {
		go_error.Throw(`amount decimal error`, constant.AMOUNT_DECIMAL_ERR)
	}

	// 校验用户最大提现金额
	if userCurrencyModel.MaxWithdrawAmount != -1 && go_decimal.Decimal.Start(params.Amount).Gt(userCurrencyModel.MaxWithdrawAmount) {
		go_error.Throw(`amount must lte max withdraw amount`, constant.USER_MAX_WITHDRAW_AMOUNT)
	}

	// 检查余额
	balance := model.BalanceLogModel.GetBalanceByUserIdCurrencyId(apiSession.UserId, currencyModel.Id)
	if go_decimal.Decimal.Start(balance.Avail).Sub(balance.Freeze).Lt(params.Amount) {
		go_error.Throw(`balance not enough`, constant.BALANCE_NOT_ENOUGH)
	}

	// 检查用户此币种每日限额
	sum := model.WithdrawModel.GetWithdrewTotalOfToday(apiSession.UserId, params.Currency, params.Chain)
	if userCurrencyModel.WithdrawLimitDaily != -1 && go_decimal.Decimal.Start(sum).Add(params.Amount).Gt(userCurrencyModel.WithdrawLimitDaily) {
		go_error.Throw(`must lt withdraw limit daily`, constant.WITHDRAW_LIMIT_DAILY)
	}

	// 校验目标地址格式是否正确
	memo := ``
	if params.Memo != nil {
		memo = *params.Memo
	}
	external_service.DepositAddressService.ValidateAddress(currencyModel.Series, params.Address, memo)

	// 有tag的话，校验tag最大长度
	if memo != `` && currencyModel.HasTag == 1 && len(memo) > int(currencyModel.MaxTagLength) {
		go_error.Throw(`memo is too long`, constant.MEMO_TOO_LONG)
	}

	// 提现二次确认
	userModel := model.TeamModel.GetByUserIdIsBanned(apiSession.UserId, false)
	if userModel == nil {
		go_error.Throw(`invalid user`, constant.ILLEGAL_USER)
	}
	// 对提现二次确认请求签名
	timestamp := go_reflect.Reflect.MustToString(time.Now().UnixNano() / 1e6)
	responseKeyModel := model.ResponseKeyModel.GetByUserId(apiSession.UserId)
	if responseKeyModel == nil {
		go_error.ThrowInternal(` - user do not have response keys.`)
	}
	go_logger.Logger.Debug(`params: `, params)
	paramsMap := map[string]interface{}{
		"address":    params.Address,
		"amount":     params.Amount,
		"chain":      params.Chain,
		"currency":   params.Currency,
		"memo":       memo,
		"request_id": params.RequestId}
	content := go_json.Json.MustStringify(paramsMap)
	sig := signature.SignMessage(content+`|`+timestamp, responseKeyModel.PrivateKey)
	go_logger.Logger.DebugF("content: %s\n", content)
	httpUtil := go_http.NewHttpRequester(go_http.WithTimeout(5 * time.Second))
	strResult := httpUtil.MustPostForString(go_http.RequestParam{
		Url:    userModel.WithdrawConfirmUrl,
		Params: params,
		Headers: map[string]interface{}{
			`STM-RES-SIGNATURE`: sig,
			`STM-RES-TIMESTAMP`: timestamp,
		},
	})
	mark := `ok`
	confirmStatus := 2
	if strResult != `ok` {
		mark = go_json.Json.MustStringify(strResult)
		confirmStatus = 3
	}
	go_mysql.MysqlHelper.MustRawExec(
		`insert into push_log (user_id,error_count,url,status,type,data, mark) values (?,?,?,?,?,?,?)`,
		apiSession.UserId,
		0,
		userModel.WithdrawConfirmUrl,
		confirmStatus,
		3,
		content,
		mark)
	if strResult != `ok` {
		go_error.Throw(`withdraw confirm failed`, constant.WITHDRAW_CONFIRM_FAIL)
	}
	tran := go_mysql.MysqlHelper.MustBegin()
	defer func() {
		if err := recover(); err != nil {
			tran.MustRollback()
			panic(err)
		} else {
			tran.MustCommit()
		}
	}()
	// 判断是否进入审核
	var status uint64
	var id uint64
	if go_decimal.Decimal.Start(params.Amount).Lte(userCurrencyModel.WithdrawCheckLimit) || go_decimal.Decimal.Start(userCurrencyModel.WithdrawCheckLimit).Eq(-1) {
		// 直接通过
		status = 3
		id = model.WithdrawModel.Insert(tran, params.RequestId, apiSession.UserId, currencyModel.Id, params.Currency, params.Chain, params.Amount, status, params.Address, memo)
	} else {
		status = 2
		id = model.WithdrawModel.Insert(tran, params.RequestId, apiSession.UserId, currencyModel.Id, params.Currency, params.Chain, params.Amount, status, params.Address, memo)
	}
	// 冻结资产
	model.BalanceLogModel.Freeze(tran, apiSession.UserId, currencyModel.Id, params.Amount, 1, id)

	return ``
}
