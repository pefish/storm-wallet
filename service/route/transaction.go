package route

import (
	"github.com/pefish/go-core/api-channel-builder"
	api_strategy2 "github.com/pefish/go-core/api-strategy"
	"wallet-storm-wallet/api-strategy"
	"wallet-storm-wallet/controller"
)

var TransactionRoute = map[string]*api_channel_builder.Route{
	`get_deposit_transaction`: {
		Description: "获取充值交易详情",
		Path:        "/v1/deposit/transaction",
		Method:      "GET",
		Strategies: []api_channel_builder.StrategyRoute{
			{
				Strategy: &api_strategy.ApikeyAuthStrategy,
				Disable:  false,
			},
		},
		ParamType: api_strategy2.ALL_TYPE,
		Params: controller.GetDepositTransactionParam{
			TxId:  `0xfeaef9adaa8a949f474cd86d573838e370441616c53836f4a19de0db64b73a68`,
		},
		Controller: controller.TransactionController.GetDepositTransaction,
		Return: api_channel_builder.ApiResult{
			Data: []controller.GetDepositTransactionReturn{
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
	`get_withdraw_transaction`: {
		Description: "获取提现交易详情",
		Path:        "/v1/withdraw/transaction",
		Method:      "GET",
		Strategies: []api_channel_builder.StrategyRoute{
			{
				Strategy: &api_strategy.ApikeyAuthStrategy,
				Disable:  false,
			},
		},
		ParamType: api_strategy2.ALL_TYPE,
		Params: controller.GetWithdrawTransactionParam{
			TxId:  `0xfeaef9adaa8a949f474cd86d573838e370441616c53836f4a19de0db64b73a68`,
		},
		Controller: controller.TransactionController.GetWithdrawTransaction,
		Return: api_channel_builder.ApiResult{
			Data: []controller.GetWithdrawTransactionReturn{
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
					OutputIndex:   0,
					Tag:           ``,
					NetworkFee:    0.0043,
					CreatedAt:     `2019-09-04T06:41:39Z`,
				},
			},
		},
	},
}
