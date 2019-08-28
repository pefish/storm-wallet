package route

import (
	"github.com/pefish/go-core/api-channel-builder"
	api_strategy2 "github.com/pefish/go-core/api-strategy"
	"wallet-storm-wallet/api-strategy"
	"wallet-storm-wallet/controller"
)

var UserRoute = map[string]*api_channel_builder.Route{
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
		ParamType:  api_strategy2.ALL_TYPE,
		Params:     controller.NewAddressParam{},
		Controller: controller.UserController.NewAddress,
	},
}
