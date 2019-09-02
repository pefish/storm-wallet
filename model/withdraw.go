package model

var WithdrawModel = Withdraw{}

type Withdraw struct {
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
	TxHex         string  `db:"tx_hex" json:"tx_hex"`
	Confirmations int64   `db:"confirmations" json:"confirmations"`
	OutputIndex   int64   `db:"output_index" json:"output_index"`
	Nonce         int64   `db:"nonce" json:"nonce"`
	NetworkFee    float64 `db:"network_fee" json:"network_fee"`
	Tag           string  `db:"tag" json:"tag"`
	Mark          string  `db:"mark" json:"mark"`

	BaseModel
}

func (this *Withdraw) GetTableName() string {
	return `withdraw`
}
