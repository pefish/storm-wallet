package controller

import (
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-mysql"
	"wallet-storm-wallet/model"
)

type UserControllerClass struct {
}

var UserController = UserControllerClass{}

type ListBalanceReturn struct {
	Currency string `json:"currency"`
	Chain    string `json:"chain"`
	Avail    string `json:"avail"`
	Freeze   string `json:"freeze"`
}

func (this *UserControllerClass) ListBalance(apiSession *api_session.ApiSessionClass) interface{} {
	results := []ListBalanceReturn{}
	model.BalanceLogModel.ListAllBalanceByUserIdForStruct(&results, apiSession.UserId)
	return results
}

func (this *UserControllerClass) GetCoinBalance(apiSession *api_session.ApiSessionClass) interface{} {
	params := GetUserCurrencyParam{}
	apiSession.ScanParams(&params)
	currencyInfo := model.CurrencyModel.GetIdByCurrencyChain(params.Currency, params.Chain)

	result := model.BalanceLogModel.GetBalanceByUserIdCurrencyId(apiSession.UserId, currencyInfo.Id)
	return result
}

type ListUserCurrencyReturn struct {
	WithdrawLimitDaily            float64 `json:"withdraw_limit_daily"`
	MaxWithdrawAmount             float64 `json:"max_withdraw_amount"`
	WithdrawCheckLimit            float64 `json:"withdraw_check_limit"`
	Currency                      string  `json:"currency"`
	Chain                         string  `json:"chain"`
	ContractAddress               string  `json:"contract_address"`
	Decimals                      uint64  `json:"decimals"`
	DepositConfirmationThreshold  uint64  `json:"deposit_confirmation_threshold"`
	WithdrawConfirmationThreshold uint64  `json:"withdraw_confirmation_threshold"`
	NetworkFeeCurrency            string  `json:"network_fee_currency"`
	NetworkFeeDecimal             uint64  `json:"network_fee_decimal"`
	HasTag                        uint64  `json:"has_tag"`
	MaxTagLength                  uint64  `json:"max_tag_length"`
	IsWithdrawEnable              uint64  `json:"is_withdraw_enable"`
	IsDepositEnable               uint64  `json:"is_deposit_enable"`
}

func (this *UserControllerClass) ListUserCurrencies(apiSession *api_session.ApiSessionClass) interface{} {
	var results []ListUserCurrencyReturn
	go_mysql.MysqlHelper.MustRawSelect(&results, `
select
a.withdraw_limit_daily,a.max_withdraw_amount,a.withdraw_check_limit,
b.currency,b.chain,b.contract_address,b.decimals,b.deposit_confirmation_threshold,b.withdraw_confirmation_threshold,
b.network_fee_currency,b.network_fee_decimal,b.has_tag,b.max_tag_length,b.is_withdraw_enable,b.is_deposit_enable
from user_currency a
left join currency b
on a.currency_id = b.id
where a.user_id = ?
`, apiSession.UserId)
	return results
}

type GetUserCurrencyParam struct {
	Currency string `json:"currency" validate:"required"`
	Chain    string `json:"chain" validate:"required"`
}

func (this *UserControllerClass) GetUserCurrency(apiSession *api_session.ApiSessionClass) interface{} {
	var param GetUserCurrencyParam
	apiSession.ScanParams(&param)
	var result ListUserCurrencyReturn
	go_mysql.MysqlHelper.MustRawSelectFirst(&result, `
select
a.withdraw_limit_daily,a.max_withdraw_amount,a.withdraw_check_limit,
b.currency,b.chain,b.contract_address,b.decimals,b.deposit_confirmation_threshold,b.withdraw_confirmation_threshold,
b.network_fee_currency,b.network_fee_decimal,b.has_tag,b.max_tag_length,b.is_withdraw_enable,b.is_deposit_enable
from user_currency a
left join currency b
on a.currency_id = b.id
where a.user_id = ? and b.currency = ? and b.chain = ?
`, apiSession.UserId, param.Currency, param.Chain)
	return result
}
