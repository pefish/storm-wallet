package controller

import (
	_type "github.com/pefish/go-core/api-session/type"
	go_error "github.com/pefish/go-error"
	"wallet-storm-wallet/model"
)

type TransactionControllerClass struct {
}

var TransactionController = TransactionControllerClass{}

type ListDepositTransactionParam struct {
	Chain *string `json:"chain,omitempty" validate:"omitempty" desc:"要查询哪条链上的交易"`
	TxId  string  `json:"tx_id" validate:"required" desc:"要查询的tx id"`
}

type ListDepositTransactionReturn struct {
	UserId        uint64  `json:"user_id"`
	Currency      string  `json:"currency"`
	Chain         string  `json:"chain"`
	Amount        float64 `json:"amount"`
	Address       string  `json:"address"`
	Status        int64   `json:"status"`
	Height        string  `json:"height"`
	BlockId       string  `json:"block_id"`
	TxId          string  `json:"tx_id"`
	Confirmations int64   `json:"confirmations"`
	OutputIndex   int64   `json:"output_index"`
	Tag           string  `json:"tag"`
	CreatedAt     string  `json:"created_at"`
}

func (this *TransactionControllerClass) ListDepositTransaction(apiSession _type.IApiSession) (interface{}, *go_error.ErrorInfo) {
	params := ListDepositTransactionParam{}
	apiSession.ScanParams(&params)

	results := []ListDepositTransactionReturn{}
	if params.Chain == nil {
		model.DepositModel.ListByUserIdTxIdForStruct(&results, apiSession.UserId(), params.TxId)
	} else {
		model.DepositModel.ListByUserIdChainTxIdForStruct(&results, apiSession.UserId(), *params.Chain, params.TxId)
	}

	return results, nil
}

type ListWithdrawTransactionParam struct {
	Chain *string `json:"chain,omitempty" validate:"omitempty" desc:"要查询哪条链上的交易"`
	TxId  string  `json:"tx_id" validate:"required" desc:"要查询的tx id"`
}

type ListWithdrawTransactionReturn struct {
	UserId        uint64  `json:"user_id"`
	Currency      string  `json:"currency"`
	Chain         string  `json:"chain"`
	Amount        float64 `json:"amount"`
	FromAddress   string  `json:"from_address"`
	ToAddress     string  `json:"to_address"`
	Status        int64   `json:"status"`
	Height        string  `json:"height"`
	BlockId       string  `json:"block_id"`
	TxId          string  `json:"tx_id"`
	Confirmations int64   `json:"confirmations"`
	NetworkFee    float64 `json:"network_fee"`
	Tag           string  `json:"tag"`
	CreatedAt     string  `json:"created_at"`
}

func (this *TransactionControllerClass) ListWithdrawTransaction(apiSession _type.IApiSession) (interface{}, *go_error.ErrorInfo) {
	params := ListWithdrawTransactionParam{}
	apiSession.ScanParams(&params)

	results := []ListWithdrawTransactionReturn{}
	if params.Chain == nil {
		model.WithdrawModel.ListByUserIdTxIdForStruct(&results, apiSession.UserId(), params.TxId)
	} else {
		model.WithdrawModel.ListByUserIdChainTxIdForStruct(&results, apiSession.UserId(), *params.Chain, params.TxId)
	}

	return results, nil
}
