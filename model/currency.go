package model

var CurrencyModel = Currency{}

type Currency struct {
	Currency             string  `db:"currency" json:"currency"`
	Chain                string  `db:"chain" json:"chain"`
	ContractAddress      string  `db:"contract_address" json:"contract_address"`
	Decimals             int64   `db:"decimals" json:"decimals"`
	MinDepositAmount     float64 `db:"min_deposit_amount" json:"min_deposit_amount"`
	MinWithdrawAmount    float64 `db:"min_withdraw_amount" json:"min_withdraw_amount"`
	MaxWithdrawAmount    float64 `db:"max_withdraw_amount" json:"max_withdraw_amount"`
	ScanLimit            float64 `db:"scan_limit" json:"scan_limit"`
	ReservedValue        float64 `db:"reserved_value" json:"reserved_value"`
	DepositConfirmation  int64   `db:"deposit_confirmation" json:"deposit_confirmation"`
	WithdrawConfirmation int64   `db:"withdraw_confirmation" json:"withdraw_confirmation"`
	Series               string  `db:"series" json:"series"`
	WithdrawInoutNum     int64   `db:"withdraw_inout_num" json:"withdraw_inout_num"`
	NetworkFee           float64 `db:"network_fee" json:"network_fee"`
	WithdrawCheckLimit   float64 `db:"withdraw_check_limit" json:"withdraw_check_limit"`
	ScanInoutNum         int64   `db:"scan_inout_num" json:"scan_inout_num"`
	HasTag               int64   `db:"has_tag" json:"has_tag"`
	DepositEnable        int64   `db:"deposit_enable" json:"deposit_enable"`
	WithdrawEnable       int64   `db:"withdraw_enable" json:"withdraw_enable"`
	WithdrawCheckAuto    int64   `db:"withdraw_check_auto" json:"withdraw_check_auto"`
	IsBanned             int64   `db:"is_banned" json:"is_banned"`
	MaxTagLength         int64   `db:"max_tag_length" json:"max_tag_length"`

	BaseModel
}

func (this *Currency) GetTableName() string {
	return `currency`
}
