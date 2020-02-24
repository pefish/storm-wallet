package route

import (
	"github.com/pefish/go-core/api"
	api_strategy2 "github.com/pefish/go-core/api-strategy"
	global_api_strategy "github.com/pefish/go-core/driver/global-api-strategy"
	"wallet-storm-wallet/api-strategy"
	"wallet-storm-wallet/controller"
	"wallet-storm-wallet/return-hook"
)

var TransactionRoute = []*api.Api{
	{
		Description: "获取充值交易详情",
		Path:        "/v1/deposit/transactions",
		Method:      "GET",
		Strategies: []global_api_strategy.StrategyData{
			{
				Strategy: &api_strategy.ApikeyAuthStrategy,
				Disable:  false,
			},
		},
		ParamType: api_strategy2.ALL_TYPE,
		Params: controller.ListDepositTransactionParam{
			TxId: `0xfeaef9adaa8a949f474cd86d573838e370441616c53836f4a19de0db64b73a68`,
		},
		Controller:     controller.TransactionController.ListDepositTransaction,
		ReturnHookFunc: return_hook.ReturnHook,
		Return: api.ApiResult{
			Data: []controller.ListDepositTransactionReturn{
				{
					UserId:        1,
					Currency:      `ETH`,
					Chain:         `Eth`,
					Amount:        4.4646,
					Address:       `0x3a7f6e30d48c9a0120926b3bc930fe1992a4592c`,
					Status:        1,
					Height:        ``,
					BlockId:       ``,
					TxId:          `0xfeaef9adaa8a949f474cd86d573838e370441616c53836f4a19de0db64b73a68`,
					Confirmations: 1,
					OutputIndex:   0,
					Tag:           ``,
					CreatedAt:     `2019-09-04T06:41:39Z`,
				},
			},
		},
	},
	{
		Description: "获取提现交易详情",
		Path:        "/v1/withdraw/transactions",
		Method:      "GET",
		Strategies: []global_api_strategy.StrategyData{
			{
				Strategy: &api_strategy.ApikeyAuthStrategy,
				Disable:  false,
			},
		},
		ParamType: api_strategy2.ALL_TYPE,
		Params: controller.ListWithdrawTransactionParam{
			TxId: `0xfeaef9adaa8a949f474cd86d573838e370441616c53836f4a19de0db64b73a68`,
		},
		Controller:     controller.TransactionController.ListWithdrawTransaction,
		ReturnHookFunc: return_hook.ReturnHook,
		Return: api.ApiResult{
			Data: []controller.ListWithdrawTransactionReturn{
				{
					UserId:        1,
					Currency:      `ETH`,
					Chain:         `Eth`,
					Amount:        4.4646,
					FromAddress:   `0x3a7f6e30d48c9a0120926b3bc930fe1992a4592c`,
					ToAddress:     `0x5cee9037344f57dbbef7de90348c5c82a1472882`,
					Status:        1,
					Height:        ``,
					BlockId:       ``,
					TxId:          `0xfeaef9adaa8a949f474cd86d573838e370441616c53836f4a19de0db64b73a68`,
					Confirmations: 1,
					Tag:           ``,
					NetworkFee:    0.0043,
					CreatedAt:     `2019-09-04T06:41:39Z`,
				},
			},
		},
	},
}
