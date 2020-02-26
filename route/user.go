package route

import (
	"github.com/pefish/go-core/api"
	api_strategy2 "github.com/pefish/go-core/api-strategy"
	global_api_strategy2 "github.com/pefish/go-core/global-api-strategy"
	"wallet-storm-wallet/api-strategy"
	"wallet-storm-wallet/controller"
	"wallet-storm-wallet/return-hook"
)

var UserRoute = []*api.Api{
	{
		Description: "获取账户余额",
		Path:        "/v1/balance",
		Method:      "GET",
		Strategies: []api_strategy2.StrategyData{
			{
				Strategy: &api_strategy.ApikeyAuthStrategy,
				Disable:  false,
			},
		},
		ParamType:      global_api_strategy2.ALL_TYPE,
		Controller:     controller.UserController.ListBalance,
		ReturnHookFunc: return_hook.ReturnHook,
		Return: api.ApiResult{
			Data: []controller.ListBalanceReturn{
				{
					Currency: `ETH`,
					Chain:    `Eth`,
					Avail:    `73567`,
					Freeze:   `77`,
				},
			},
		},
	},
	{
		Description: "获取用户开启的所有币种",
		Path:        "/v1/user-currencies",
		Method:      "GET",
		Strategies: []api_strategy2.StrategyData{
			{
				Strategy: &api_strategy.ApikeyAuthStrategy,
				Disable:  false,
			},
		},
		ParamType:      global_api_strategy2.ALL_TYPE,
		Controller:     controller.UserController.ListUserCurrencies,
		ReturnHookFunc: return_hook.ReturnHook,
		Return: api.ApiResult{
			Data: []controller.ListUserCurrencyReturn{
				{
					WithdrawLimitDaily:            200,
					MaxWithdrawAmount:             100,
					WithdrawCheckLimit:            100,
					Currency:                      `ETH`,
					Chain:                         `Eth`,
					ContractAddress:               ``,
					Decimals:                      18,
					DepositConfirmationThreshold:  12,
					WithdrawConfirmationThreshold: 30,
					NetworkFeeCurrency:            `Eth.ETH`,
					NetworkFeeDecimal:             18,
					HasTag:                        0,
					MaxTagLength:                  0,
					IsWithdrawEnable:              1,
					IsDepositEnable:               1,
				},
			},
		},
	},
	{
		Description: "获取用户开启的币种",
		Path:        "/v1/user-currency",
		Method:      "GET",
		Strategies: []api_strategy2.StrategyData{
			{
				Strategy: &api_strategy.ApikeyAuthStrategy,
				Disable:  false,
			},
		},
		ParamType: global_api_strategy2.ALL_TYPE,
		Params: controller.GetUserCurrencyParam{
			Currency: `ETH`,
			Chain:    `Eth`,
		},
		Controller:     controller.UserController.GetUserCurrency,
		ReturnHookFunc: return_hook.ReturnHook,
		Return: api.ApiResult{
			Data: []controller.ListUserCurrencyReturn{
				{
					WithdrawLimitDaily:            200,
					MaxWithdrawAmount:             100,
					WithdrawCheckLimit:            100,
					Currency:                      `ETH`,
					Chain:                         `Eth`,
					ContractAddress:               ``,
					Decimals:                      18,
					DepositConfirmationThreshold:  12,
					WithdrawConfirmationThreshold: 30,
					NetworkFeeCurrency:            `Eth.ETH`,
					NetworkFeeDecimal:             18,
					HasTag:                        0,
					MaxTagLength:                  0,
					IsWithdrawEnable:              1,
					IsDepositEnable:               1,
				},
			},
		},
	},
	{
		Description: "获取账户指定币种余额",
		Path:        "/v1/coin-balance",
		Method:      "GET",
		Strategies: []api_strategy2.StrategyData{
			{
				Strategy: &api_strategy.ApikeyAuthStrategy,
				Disable:  false,
			},
		},
		ParamType:      global_api_strategy2.ALL_TYPE,
		Controller:     controller.UserController.GetCoinBalance,
		ReturnHookFunc: return_hook.ReturnHook,
		Return: api.ApiResult{
			Data: []controller.ListBalanceReturn{
				{
					Currency: `ETH`,
					Chain:    `Eth`,
					Avail:    `73`,
					Freeze:   `7`,
				},
			},
		},
	},
}
