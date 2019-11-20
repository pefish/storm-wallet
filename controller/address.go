package controller

import (
	"fmt"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-error"
	"github.com/pefish/go-redis"
	"github.com/satori/go.uuid"
	"time"
	"wallet-storm-wallet/constant"
	"wallet-storm-wallet/external-service"
	"wallet-storm-wallet/model"
)

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
	params := NewAddressParam{}
	apiSession.ScanParams(&params)

	lockKey := fmt.Sprintf(`storm_wallet_get_deposit_address_%d_lock`, apiSession.UserId)
	uniqueId := uuid.NewV1().String()
	if !go_redis.RedisHelper.GetLock(lockKey, uniqueId, 3*time.Second) {
		go_error.Throw(`rate limit`, constant.API_RATELIMIT)
	}
	defer go_redis.RedisHelper.ReleaseLock(lockKey, uniqueId)

	currencyModel := model.UserCurrencyModel.GetCurrencyOfUserByCurrencyChain(apiSession.UserId, params.Currency, params.Chain)
	if currencyModel == nil {
		go_error.Throw(`user currency is not available`, constant.USER_CURRENCY_NOT_AVAILABLE)
	}
	depositAddressModel := model.DepositAddressModel.GetByUserIdSeriesIndex(apiSession.UserId, currencyModel.Series, params.Index)
	if depositAddressModel != nil {
		return NewAddressReturn{
			Address: depositAddressModel.Address,
		}
	}

	result := external_service.DepositAddressService.GetAddress(currencyModel.Series, apiSession.UserId, params.Index)
	model.DepositAddressModel.Insert(apiSession.UserId, result.Address, result.Path, currencyModel.Series, params.Index)
	return NewAddressReturn{
		Address: result.Address,
	}
}

type ValidateAddressParam struct {
	Currency string `json:"currency" validate:"required" desc:"currency"`
	Chain    string `json:"chain" validate:"required" desc:"要验证哪条链上的地址"`
	Address  string `json:"address" validate:"required" desc:"address"`
}

func (this *AddressControllerClass) ValidateAddress(apiSession *api_session.ApiSessionClass) interface{} {
	params := ValidateAddressParam{}
	apiSession.ScanParams(&params)

	currencyModel := model.UserCurrencyModel.GetCurrencyOfUserByCurrencyChain(apiSession.UserId, params.Currency, params.Chain)
	if currencyModel == nil {
		go_error.Throw(`user currency is not available`, constant.USER_CURRENCY_NOT_AVAILABLE)
	}
	external_service.DepositAddressService.ValidateAddress(currencyModel.Series, params.Address, ``)
	return true
}

type IsPlatformAddressParam struct {
	Currency string `json:"currency" validate:"required" desc:"currency"`
	Chain    string `json:"chain" validate:"required" desc:"要查询哪条链上的地址"`
	Address  string `json:"address" validate:"required" desc:"address"`
}

func (this *AddressControllerClass) IsPlatformAddress(apiSession *api_session.ApiSessionClass) interface{} {
	params := IsPlatformAddressParam{}
	apiSession.ScanParams(&params)

	currencyModel := model.UserCurrencyModel.GetCurrencyOfUserByCurrencyChain(apiSession.UserId, params.Currency, params.Chain)
	if currencyModel == nil {
		go_error.Throw(`user currency is not available`, constant.USER_CURRENCY_NOT_AVAILABLE)
	}
	depositAddressModel := model.DepositAddressModel.GetByUserIdSeriesAddress(apiSession.UserId, currencyModel.Series, params.Address)
	return depositAddressModel != nil
}
