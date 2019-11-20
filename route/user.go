package route

import (
	"github.com/pefish/go-core/api-channel-builder"
	api_strategy2 "github.com/pefish/go-core/api-strategy"
	"wallet-storm-wallet/api-strategy"
	"wallet-storm-wallet/controller"
)

var UserRoute = map[string]*api_channel_builder.Route{
	`get_balance`: {
		Description: "获取账户余额",
		Path:        "/v1/balance",
		Method:      "GET",
		Strategies: []api_channel_builder.StrategyRoute{
			{
				Strategy: &api_strategy.ApikeyAuthStrategy,
				Disable:  false,
			},
		},
		ParamType:  api_strategy2.ALL_TYPE,
		Controller: controller.UserController.ListBalance,
		Return: api_channel_builder.ApiResult{
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
	//`get_user_currencies`: {
	//	Description: "获取账户所有币种",
	//	Path:        "/v1/user-currency",
	//	Method:      "GET",
	//	Strategies: []api_channel_builder.StrategyRoute{
	//		{
	//			Strategy: &api_strategy.ApikeyAuthStrategy,
	//			Disable:  false,
	//		},
	//	},
	//	ParamType:  api_strategy2.ALL_TYPE,
	//	Controller: controller.UserController.ListUserCurrency,
	//	Return: api_channel_builder.ApiResult{
	//		Data: []controller.ListUserCurrencyReturn{
	//			WithdrawLimitDaily:    `67372`,
	//			MaxWithdrawAmount:     `52544`,
	//			WithdrawCheckLimit:    `73567`,
	//			Currency:              `ETH`,
	//			Chain:                 `Eth`,
	//			ContractAddress:       `0xywtrywhsthy`,
	//			Decimals:              8,
	//			min_deposit_amount:    `12`,
	//			min_withdraw_amount:   `7365`,
	//			deposit_confirmation:  12,
	//			withdraw_confirmation: 30,
	//			network_fee:           `0.66`,
	//			has_tag:               0,
	//			max_tag_length:        150,
	//			is_deposit_enable:     0,
	//			is_withdraw_enable:    1,
	//		},
	//	},
	//},
	//`get_user_addresses`: {
	//	Description: "获取账户充值地址",
	//	Path:        "/v1/user-address",
	//	Method:      "GET",
	//	Strategies: []api_channel_builder.StrategyRoute{
	//		{
	//			Strategy: &api_strategy.ApikeyAuthStrategy,
	//			Disable:  false,
	//		},
	//	},
	//	ParamType: api_strategy2.ALL_TYPE,
	//	Params: controller.ListUserAddressParam{
	//		Page:     1,
	//		Size:     10,
	//		Chain:    `Eth`,
	//		Currency: `ETH`,
	//	},
	//	Controller: controller.UserController.ListUserAddress,
	//	Return: api_channel_builder.ApiResult{
	//		Data: UserAddressReturn{
	//			List: UserAddressListReturn{
	//				Address: `0xhdghytwthetyt33y3y`,
	//				Index:   63562,
	//			},
	//			Count: 30,
	//		},
	//	},
	//},
}
