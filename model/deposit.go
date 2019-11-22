package model

import (
	"github.com/pefish/go-mysql"
)

var DepositModel = Deposit{}

type Deposit struct {
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
	Mark          string  `json:"mark"`

	BaseModel
}

func (this *Deposit) GetTableName() string {
	return `deposit`
}

func (this *Deposit) ListByUserIdChainTxIdForStruct(results interface{}, userId uint64, chain string, txId string) {
	go_mysql.MysqlHelper.Select(results, this.GetTableName(), `*`, map[string]interface{}{
		`user_id`: userId,
		`chain`:   chain,
		`tx_id`:   txId,
	})
}

func (this *Deposit) ListByUserIdTxIdForStruct(results interface{}, userId uint64, txId string) {
	go_mysql.MysqlHelper.Select(results, this.GetTableName(), `*`, map[string]interface{}{
		`user_id`: userId,
		`tx_id`:   txId,
	})
}
