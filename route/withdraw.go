package route

import (
	"github.com/pefish/go-core/api-channel-builder"
	api_strategy2 "github.com/pefish/go-core/api-strategy"
	"wallet-storm-wallet/api-strategy"
	"wallet-storm-wallet/controller"
	"wallet-storm-wallet/return-hook"
)

var WithdrawRoute = map[string]*api_channel_builder.Route{
	`withdraw`: {
		Description: "发起提现",
		Path:        "/v1/withdraw",
		Method:      "POST",
		Strategies: []api_channel_builder.StrategyRoute{
			{
				Strategy: &api_strategy.ApikeyAuthStrategy,
				Disable:  false,
				Param: api_strategy.ApikeyAuthParam{
					AllowedType: `2`,
				},
			},
		},
		ParamType: api_strategy2.ALL_TYPE,
		Params: controller.WithdrawParam{
			Currency:  `ETH`,
			Chain:     `Eth`,
			RequestId: `hsgfh65`,
			Address:   `0xfb6d58f5dc77ff06390fe1f30c57e670b555b34a`,
			Amount:    `0.6`,
		},
		Controller:     controller.WithdrawController.Withdraw,
		ReturnHookFunc: return_hook.ReturnHook,
		Return: api_channel_builder.ApiResult{
			Data: ``,
		},
	},
}
