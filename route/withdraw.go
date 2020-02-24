package route

import (
	"github.com/pefish/go-core/api"
	api_strategy2 "github.com/pefish/go-core/api-strategy"
	global_api_strategy "github.com/pefish/go-core/driver/global-api-strategy"
	"wallet-storm-wallet/api-strategy"
	"wallet-storm-wallet/controller"
	"wallet-storm-wallet/return-hook"
)

var WithdrawRoute = []*api.Api{
	{
		Description: "发起提现",
		Path:        "/v1/withdraw",
		Method:      "POST",
		Strategies: []global_api_strategy.StrategyData{
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
		Return: api.ApiResult{
			Data: ``,
		},
	},
}
