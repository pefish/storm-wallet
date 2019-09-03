package controller

import "github.com/pefish/go-core/api-session"

type WithdrawControllerClass struct {
}

var WithdrawController = WithdrawControllerClass{}

type WithdrawParam struct {
	Currency  string `json:"currency" validate:"required" desc:"提现的币种名"`
	Chain     string `json:"chain" validate:"required" desc:"提现的链名"`
	RequestId string `json:"request_id" validate:"required" desc:"订单id。此id幂等"`
	Address   string `json:"address" validate:"required" desc:"提现的目标地址"`
	Amount    string `json:"amount" validate:"required" desc:"提现的数量"`
	Memo      string `json:"memo" validate:"omitempty" desc:"memo"`
}

func (this *WithdrawControllerClass) Withdraw(apiSession *api_session.ApiSessionClass) interface{} {
	return ``
}
