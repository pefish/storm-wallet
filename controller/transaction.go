package controller

import (
	"github.com/pefish/go-core/api-session"
	"wallet-storm-wallet/model"
)

type TransactionControllerClass struct {
}

var TransactionController = TransactionControllerClass{}

type GetDepositTransactionParam struct {
	Chain *string `json:"chain" validate:"omitempty" desc:"要查询哪条链上的交易"`
	TxId  string  `json:"tx_id" validate:"required" desc:"要查询的tx id"`
}

type GetDepositTransactionReturn struct {
	UserId        uint64  `db:"user_id" json:"user_id"`
	Currency      string  `db:"currency" json:"currency"`
	Chain         string  `db:"chain" json:"chain"`
	Amount        float64 `db:"amount" json:"amount"`
	Address       string  `db:"address" json:"address"`
	Status        int64   `db:"status" json:"status"`
	Height        string  `db:"height" json:"height"`
	BlockId       string  `db:"block_id" json:"block_id"`
	TxId          string  `db:"tx_id" json:"tx_id"`
	Confirmations int64   `db:"confirmations" json:"confirmations"`
	OutputIndex   int64   `db:"output_index" json:"output_index"`
	Tag           string  `db:"tag" json:"tag"`
	CreatedAt     string  `db:"created_at" json:"created_at"`
}

func (this *TransactionControllerClass) GetDepositTransaction(apiSession *api_session.ApiSessionClass) interface{} {
	params := GetDepositTransactionParam{}
	apiSession.ScanParams(&params)

	results := []GetDepositTransactionReturn{}
	if params.Chain == nil {
		model.DepositModel.ListByUserIdTxIdForStruct(&results, apiSession.UserId, params.TxId)
	} else {
		model.DepositModel.ListByUserIdChainTxIdForStruct(&results, apiSession.UserId, *params.Chain, params.TxId)
	}

	return results
}

type GetWithdrawTransactionParam struct {
	Chain *string `json:"chain" validate:"omitempty" desc:"要查询哪条链上的交易"`
	TxId  string  `json:"tx_id" validate:"required" desc:"要查询的tx id"`
}

type GetWithdrawTransactionReturn struct {
	UserId        uint64  `db:"user_id" json:"user_id"`
	Currency      string  `db:"currency" json:"currency"`
	Chain         string  `db:"chain" json:"chain"`
	Amount        float64 `db:"amount" json:"amount"`
	FromAddress   string  `db:"from_address" json:"from_address"`
	ToAddress     string  `db:"to_address" json:"to_address"`
	Status        int64   `db:"status" json:"status"`
	Height        string  `db:"height" json:"height"`
	BlockId       string  `db:"block_id" json:"block_id"`
	TxId          string  `db:"tx_id" json:"tx_id"`
	Confirmations int64   `db:"confirmations" json:"confirmations"`
	OutputIndex   int64   `db:"output_index" json:"output_index"`
	NetworkFee    float64 `db:"network_fee" json:"network_fee"`
	Tag           string  `db:"tag" json:"tag"`
	CreatedAt     string  `db:"created_at" json:"created_at"`
}

func (this *TransactionControllerClass) GetWithdrawTransaction(apiSession *api_session.ApiSessionClass) interface{} {
	params := GetWithdrawTransactionParam{}
	apiSession.ScanParams(&params)

	results := []GetWithdrawTransactionReturn{}
	if params.Chain == nil {
		model.WithdrawModel.ListByUserIdTxIdForStruct(&results, apiSession.UserId, params.TxId)
	} else {
		model.WithdrawModel.ListByUserIdChainTxIdForStruct(&results, apiSession.UserId, *params.Chain, params.TxId)
	}

	return results
}
