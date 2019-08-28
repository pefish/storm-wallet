package controller

import "github.com/pefish/go-core/api-session"

type UserControllerClass struct {
}

var UserController = UserControllerClass{}

type NewAddressParam struct {
	Currency string `json:"currency" validate:"required" desc:"currency" example:"ETH,BTC"`
}
func (this *UserControllerClass) NewAddress(apiSession *api_session.ApiSessionClass) interface{} {
	return ``
}