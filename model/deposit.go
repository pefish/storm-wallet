package model

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
