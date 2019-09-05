package controller

import (
	"github.com/pefish/go-core/api-session"
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
