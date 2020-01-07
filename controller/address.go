package controller

import (
	"fmt"
	"time"
	"wallet-storm-wallet/constant"
	external_service "wallet-storm-wallet/external-service"
	"wallet-storm-wallet/model"

	api_session "github.com/pefish/go-core/api-session"
	go_crypto "github.com/pefish/go-crypto"
	go_error "github.com/pefish/go-error"
	go_redis "github.com/pefish/go-redis"
	go_reflect "github.com/pefish/go-reflect"
	uuid "github.com/satori/go.uuid"
)

type AddressControllerClass struct {
}

var AddressController = AddressControllerClass{}

type NewAddressParam struct {
	Currency string `json:"currency" validate:"required" desc:"currency"`
	Chain    string `json:"chain" validate:"required" desc:"要获取哪条链上的地址"`
	Index    uint64 `json:"index" validate:"required,max=99999999" desc:"地址索引。索引一样则返回的地址一样"`
}
type NewAddressReturn struct {
	Address string `json:"address"`
	Tag     string `json:"tag"`
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
	if currencyModel.IsDepositEnable == 0 {
		go_error.Throw(`currency deposit is not available`, constant.CURRENCY_DEPOSIT_BANNED)
	}
	depositAddressModel := model.DepositAddressModel.GetByUserIdSeriesIndex(apiSession.UserId, currencyModel.Series, params.Index)
	if depositAddressModel != nil {
		tag := ``
		if currencyModel.HasTag == 1 {
			tag = depositAddressModel.Tag
		}
		return NewAddressReturn{
			Address: depositAddressModel.Address,
			Tag:     tag,
		}
	}
	// 如果是带有tag的币种，就取热钱包地址，tag使用 位移算法（userid+index）
	if currencyModel.HasTag == 1 {
		if params.Index > 99999999 {
			go_error.Throw(`index is too big`, constant.ADDRESS_INDEX_TOO_BIG)
		}
		if apiSession.UserId > 99 { // 1-99的用户才允许取带有tag的地址
			go_error.Throw(`user id is too big`, constant.USERID_TOO_BIG)
		}
		tempTag := go_reflect.Reflect.MustToString(apiSession.UserId) + go_reflect.Reflect.MustToString(params.Index)
		tag := go_crypto.Crypto.ShiftCryptForInt(5727262753, go_reflect.Reflect.MustToInt64(tempTag))
		tagStr := go_reflect.Reflect.MustToString(tag)
		walletConfigModel := model.WalletConfigModel.GetByChainType(currencyModel.Chain, 1)
		if walletConfigModel == nil {
			go_error.Throw(`hot wallet config error`, constant.WALLET_CONFIG_ERROR)
		}
		model.DepositAddressModel.Insert(apiSession.UserId, walletConfigModel.Address, ``, currencyModel.Series, params.Index, tagStr)
		return NewAddressReturn{
			Address: walletConfigModel.Address,
			Tag:     tagStr,
		}
	}

	result := external_service.DepositAddressService.GetAddress(currencyModel.Series, apiSession.UserId, params.Index)
	if result.Address == `` {
		go_error.Throw(`address service return a null address`, constant.ILLEGAL_ADDRESS)
	}
	model.DepositAddressModel.Insert(apiSession.UserId, result.Address, result.Path, currencyModel.Series, params.Index, ``)
	return NewAddressReturn{
		Address: result.Address,
		Tag:     ``,
	}
}

type ValidateAddressParam struct {
	Currency string  `json:"currency" validate:"required" desc:"currency"`
	Chain    string  `json:"chain" validate:"required" desc:"要验证哪条链上的地址"`
	Address  string  `json:"address" validate:"required" desc:"address"`
	Tag      *string `json:"tag,omitempty" validate:"omitempty" desc:"memo"`
}

func (this *AddressControllerClass) ValidateAddress(apiSession *api_session.ApiSessionClass) interface{} {
	params := ValidateAddressParam{}
	apiSession.ScanParams(&params)

	currencyModel := model.UserCurrencyModel.GetCurrencyOfUserByCurrencyChain(apiSession.UserId, params.Currency, params.Chain)
	if currencyModel == nil {
		go_error.Throw(`user currency is not available`, constant.USER_CURRENCY_NOT_AVAILABLE)
	}
	memo := ``
	if params.Tag != nil {
		memo = *params.Tag
	}
	external_service.DepositAddressService.ValidateAddress(currencyModel.Series, params.Address, memo)
	return true
}

type IsPlatformAddressParam struct {
	Currency string  `json:"currency" validate:"required" desc:"currency"`
	Chain    string  `json:"chain" validate:"required" desc:"要查询哪条链上的地址"`
	Address  string  `json:"address" validate:"required" desc:"address"`
	Tag      *string `json:"tag,omitempty" validate:"omitempty" desc:"memo"`
}

func (this *AddressControllerClass) IsPlatformAddress(apiSession *api_session.ApiSessionClass) interface{} {
	params := IsPlatformAddressParam{}
	apiSession.ScanParams(&params)

	currencyModel := model.UserCurrencyModel.GetCurrencyOfUserByCurrencyChain(apiSession.UserId, params.Currency, params.Chain)
	if currencyModel == nil {
		go_error.Throw(`user currency is not available`, constant.USER_CURRENCY_NOT_AVAILABLE)
	}
	memo := ``
	if params.Tag != nil {
		memo = *params.Tag
	}
	depositAddressModel := model.DepositAddressModel.GetByUserIdSeriesAddress(apiSession.UserId, currencyModel.Series, params.Address, memo)
	return depositAddressModel != nil
}
