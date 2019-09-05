package model

import "github.com/pefish/go-mysql"

var DepositModel = Deposit{}

type Deposit struct {
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
	Mark          string  `db:"mark" json:"mark"`
	ScanStatus    int64   `db:"scan_status" json:"scan_status"`

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
