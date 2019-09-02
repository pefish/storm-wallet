package controller

import "github.com/pefish/go-core/api-session"

type AddressControllerClass struct {
}

var AddressController = AddressControllerClass{}

type NewAddressParam struct {
	Currency string `json:"currency" validate:"required" desc:"currency"`
	Chain    string `json:"chain" validate:"required" desc:"要获取哪条链上的地址"`
	Index    uint64 `json:"index" validate:"required,max=10000000" desc:"地址索引。索引一样则返回的地址一样"`
}
type NewAddressReturn struct {
	Address string `json:"address"`
}

func (this *AddressControllerClass) NewAddress(apiSession *api_session.ApiSessionClass) interface{} {
	return ``
}

type ValidateAddressParam struct {
	Currency string `json:"currency" validate:"required" desc:"currency"`
	Chain    string `json:"chain" validate:"required" desc:"要验证哪条链上的地址"`
	Address  string `json:"address" validate:"required" desc:"address"`
}

func (this *AddressControllerClass) ValidateAddress(apiSession *api_session.ApiSessionClass) interface{} {
	return ``
}

type IsPlatformAddressParam struct {
	Currency string `json:"currency" validate:"required" desc:"currency"`
	Chain    string `json:"chain" validate:"required" desc:"要查询哪条链上的地址"`
	Address  string `json:"address" validate:"required" desc:"address"`
}

func (this *AddressControllerClass) IsPlatformAddress(apiSession *api_session.ApiSessionClass) interface{} {
	return ``
}
