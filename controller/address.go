package controller

import (
	"errors"
	"fmt"
	_type "github.com/pefish/go-core/api-session/type"
	"time"
	"wallet-storm-wallet/constant"
	external_service "wallet-storm-wallet/external-service"
	"wallet-storm-wallet/model"

	go_error "github.com/pefish/go-error"
	go_redis "github.com/pefish/go-redis"
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

func (this *AddressControllerClass) NewAddress(apiSession _type.IApiSession) (interface{}, *go_error.ErrorInfo) {
	params := NewAddressParam{}
	apiSession.ScanParams(&params)

	lockKey := fmt.Sprintf(`storm_wallet_get_deposit_address_%d_lock`, apiSession.UserId)
	uniqueId := uuid.NewV1().String()
	if !go_redis.RedisHelper.MustGetLock(lockKey, uniqueId, 3*time.Second) {
		return nil, go_error.WrapWithAll(errors.New(`rate limit`), constant.API_RATELIMIT, nil)
	}
	defer go_redis.RedisHelper.MustReleaseLock(lockKey, uniqueId)

	currencyModel := model.UserCurrencyModel.GetCurrencyOfUserByCurrencyChain(apiSession.UserId(), params.Currency, params.Chain)
	if currencyModel == nil {
		return nil, go_error.WrapWithAll(errors.New(`user currency is not available`), constant.USER_CURRENCY_NOT_AVAILABLE, nil)
	}
	if currencyModel.IsDepositEnable == 0 {
		return nil, go_error.WrapWithAll(errors.New(`currency deposit is not available`), constant.CURRENCY_DEPOSIT_BANNED, nil)
	}
	depositAddressModel := model.DepositAddressModel.GetByUserIdSeriesIndex(apiSession.UserId(), currencyModel.Series, params.Index)
	if depositAddressModel != nil {
		return NewAddressReturn{
			Address: depositAddressModel.Address,
			Tag:     ``,
		}, nil
	}
	// 如果是带有tag的币种，就取热钱包地址, tag由客户管理
	if currencyModel.HasTag == 1 {
		if apiSession.UserId() > 99 { // 1-99的用户才允许取带有tag的地址
			return nil, go_error.WrapWithAll(errors.New(`user id is too big`), constant.USERID_TOO_BIG, nil)
		}
		// type 1 热钱包
		walletConfigModel := model.WalletConfigModel.GetByChainType(currencyModel.Chain, 1)
		if walletConfigModel == nil {
			return nil, go_error.WrapWithAll(errors.New(`hot wallet config error`), constant.WALLET_CONFIG_ERROR, nil)
		}
		model.DepositAddressModel.Insert(apiSession.UserId(), walletConfigModel.Address, ``, currencyModel.Series, params.Index, ``)
		return NewAddressReturn{
			Address: walletConfigModel.Address,
			Tag:     ``,
		}, nil
	}

	result := external_service.DepositAddressService.GetAddress(currencyModel.Series, apiSession.UserId(), params.Index)
	if result.Address == `` {
		return nil, go_error.WrapWithAll(errors.New(`address service return a null address`), constant.ILLEGAL_ADDRESS, nil)
	}
	model.DepositAddressModel.Insert(apiSession.UserId(), result.Address, result.Path, currencyModel.Series, params.Index, ``)
	return NewAddressReturn{
		Address: result.Address,
		Tag:     ``,
	}, nil
}

type ValidateAddressParam struct {
	Currency string  `json:"currency" validate:"required" desc:"currency"`
	Chain    string  `json:"chain" validate:"required" desc:"要验证哪条链上的地址"`
	Address  string  `json:"address" validate:"required" desc:"address"`
	Tag      *string `json:"tag,omitempty" validate:"omitempty" desc:"memo"`
}

func (this *AddressControllerClass) ValidateAddress(apiSession _type.IApiSession) (interface{}, *go_error.ErrorInfo) {
	params := ValidateAddressParam{}
	apiSession.ScanParams(&params)

	currencyModel := model.UserCurrencyModel.GetCurrencyOfUserByCurrencyChain(apiSession.UserId(), params.Currency, params.Chain)
	if currencyModel == nil {
		return nil, go_error.WrapWithAll(errors.New(`user currency is not available`), constant.USER_CURRENCY_NOT_AVAILABLE, nil)
	}
	memo := ``
	if params.Tag != nil {
		memo = *params.Tag
	}
	external_service.DepositAddressService.ValidateAddress(currencyModel.Series, params.Address, memo)
	return true, nil
}

type IsPlatformAddressParam struct {
	Currency string  `json:"currency" validate:"required" desc:"currency"`
	Chain    string  `json:"chain" validate:"required" desc:"要查询哪条链上的地址"`
	Address  string  `json:"address" validate:"required" desc:"address"`
	Tag      *string `json:"tag,omitempty" validate:"omitempty" desc:"memo"`
}

func (this *AddressControllerClass) IsPlatformAddress(apiSession _type.IApiSession) (interface{}, *go_error.ErrorInfo) {
	params := IsPlatformAddressParam{}
	apiSession.ScanParams(&params)

	currencyModel := model.UserCurrencyModel.GetCurrencyOfUserByCurrencyChain(apiSession.UserId(), params.Currency, params.Chain)
	if currencyModel == nil {
		return nil, go_error.WrapWithAll(errors.New(`user currency is not available`), constant.USER_CURRENCY_NOT_AVAILABLE, nil)
	}
	memo := ``
	if params.Tag != nil {
		memo = *params.Tag
	}
	depositAddressModel := model.DepositAddressModel.GetByUserIdSeriesAddress(apiSession.UserId(), currencyModel.Series, params.Address, memo)
	return depositAddressModel != nil, nil
}
