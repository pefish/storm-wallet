package route

import (
	"github.com/pefish/go-core/api"
	api_strategy2 "github.com/pefish/go-core/api-strategy"
	global_api_strategy2 "github.com/pefish/go-core/global-api-strategy"
	"wallet-storm-wallet/api-strategy"
	"wallet-storm-wallet/controller"
	"wallet-storm-wallet/return-hook"
)

var AddressRoute = []*api.Api{
	{
		Description: "获取新充值地址",
		Path:        "/v1/new-address",
		Method:      "POST",
		Strategies: []api_strategy2.StrategyData{
			{
				Strategy: &api_strategy.ApikeyAuthStrategy,
				Disable:  false,
				Param: api_strategy.ApikeyAuthParam{
					AllowedType: `2`,
				},
			},
		},
		ParamType: global_api_strategy2.ALL_TYPE,
		Params: controller.NewAddressParam{
			Currency: `ETH`,
			Chain:    `Eth`,
			Index:    1000,
		},
		Controller:     controller.AddressController.NewAddress,
		ReturnHookFunc: return_hook.ReturnHook,
		Return: api.ApiResult{
			Data: controller.NewAddressReturn{
				Address: `0xfb6d58f5dc77ff06390fe1f30c57e670b555b34a`,
			},
		},
	},
	{
		Description: "校验地址格式是否合法",
		Path:        "/v1/validate-address",
		Method:      "GET",
		Strategies: []api_strategy2.StrategyData{
			{
				Strategy: &api_strategy.ApikeyAuthStrategy,
				Disable:  false,
			},
		},
		ParamType: global_api_strategy2.ALL_TYPE,
		Params: controller.ValidateAddressParam{
			Address:  `0xfb6d58f5dc77ff06390fe1f30c57e670b555b34a`,
			Currency: `ETH`,
			Chain:    `Eth`,
		},
		Controller:     controller.AddressController.ValidateAddress,
		ReturnHookFunc: return_hook.ReturnHook,
		Return: api.ApiResult{
			Data: true,
		},
	},
	{
		Description: "校验地址是否用户平台地址",
		Path:        "/v1/is-platform-address",
		Method:      "GET",
		Strategies: []api_strategy2.StrategyData{
			{
				Strategy: &api_strategy.ApikeyAuthStrategy,
				Disable:  false,
			},
		},
		ParamType: global_api_strategy2.ALL_TYPE,
		Params: controller.IsPlatformAddressParam{
			Address:  `0xfb6d58f5dc77ff06390fe1f30c57e670b555b34a`,
			Currency: `ETH`,
			Chain:    `Eth`,
		},
		Controller:     controller.AddressController.IsPlatformAddress,
		ReturnHookFunc: return_hook.ReturnHook,
		Return: api.ApiResult{
			Data: true,
		},
	},
}
