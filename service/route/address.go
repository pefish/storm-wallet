package route

import (
	"github.com/pefish/go-core/api-channel-builder"
	api_strategy2 "github.com/pefish/go-core/api-strategy"
	"wallet-storm-wallet/api-strategy"
	"wallet-storm-wallet/controller"
)

var AddressRoute = map[string]*api_channel_builder.Route{
	`get_new_deposit_address`: {
		Description: "获取新充值地址",
		Path:        "/v1/new-address",
		Method:      "POST",
		Strategies: []api_channel_builder.StrategyRoute{
			{
				Strategy: &api_strategy.ApikeyAuthStrategy,
				Disable:  false,
			},
		},
		ParamType: api_strategy2.ALL_TYPE,
		Params: controller.NewAddressParam{
			Currency: `ETH`,
			Chain:    `Eth`,
			Index:    1000,
		},
		Controller: controller.AddressController.NewAddress,
		Return: api_channel_builder.ApiResult{
			Data: controller.NewAddressReturn{
				Address: `0xfb6d58f5dc77ff06390fe1f30c57e670b555b34a`,
			},
		},
	},
	`validate_address`: {
		Description: "校验地址格式是否合法",
		Path:        "/v1/validate-address",
		Method:      "GET",
		Strategies: []api_channel_builder.StrategyRoute{
			{
				Strategy: &api_strategy.ApikeyAuthStrategy,
				Disable:  false,
			},
		},
		ParamType: api_strategy2.ALL_TYPE,
		Params: controller.ValidateAddressParam{
			Address:  `0xfb6d58f5dc77ff06390fe1f30c57e670b555b34a`,
			Currency: `ETH`,
			Chain:    `Eth`,
		},
		Controller: controller.AddressController.ValidateAddress,
		Return: api_channel_builder.ApiResult{
			Data: true,
		},
	},
	`is_platform_address`: {
		Description: "校验地址是否平台地址",
		Path:        "/v1/is-platform-address",
		Method:      "GET",
		Strategies: []api_channel_builder.StrategyRoute{
			{
				Strategy: &api_strategy.ApikeyAuthStrategy,
				Disable:  false,
			},
		},
		ParamType: api_strategy2.ALL_TYPE,
		Params: controller.IsPlatformAddressParam{
			Address:  `0xfb6d58f5dc77ff06390fe1f30c57e670b555b34a`,
			Currency: `ETH`,
			Chain:    `Eth`,
		},
		Controller: controller.AddressController.IsPlatformAddress,
		Return: api_channel_builder.ApiResult{
			Data: true,
		},
	},
}
