package model

import "github.com/pefish/go-mysql"

var CurrencyModel = Currency{}

type Currency struct {
	Currency                  string  `db:"currency" json:"currency"`
	Chain                     string  `db:"chain" json:"chain"`
	ContractAddress           string  `db:"contract_address" json:"contract_address"`
	Decimals                  int64   `db:"decimals" json:"decimals"`
	MinDepositAmount          float64 `db:"min_deposit_amount" json:"min_deposit_amount"`
	MinWithdrawAmount         float64 `db:"min_withdraw_amount" json:"min_withdraw_amount"`
	ScanLimit                 float64 `db:"scan_limit" json:"scan_limit"`
	ReservedValue             float64 `db:"reserved_value" json:"reserved_value"`
	DepositConfirmation       int64   `db:"deposit_confirmation" json:"deposit_confirmation"`
	WithdrawConfirmation      int64   `db:"withdraw_confirmation" json:"withdraw_confirmation"`
	Series                    string  `db:"series" json:"series"`
	WithdrawInoutNum          int64   `db:"withdraw_inout_num" json:"withdraw_inout_num"`
	WithdrawNetworkFee        float64 `db:"withdraw_network_fee" json:"withdraw_network_fee"`
	ScanNetworkFee            float64 `db:"scan_network_fee" json:"scan_network_fee"`
	ScanInoutNum              int64   `db:"scan_inout_num" json:"scan_inout_num"`
	HasTag                    int64   `db:"has_tag" json:"has_tag"`
	IsWithdrawEnable          int64   `db:"is_withdraw_enable" json:"is_withdraw_enable"`
	IsBanned                  int64   `db:"is_banned" json:"is_banned"`
	MaxTagLength              int64   `db:"max_tag_length" json:"max_tag_length"`
	DefaultWithdrawLimitDaily float64 `db:"default_withdrawlimit_daily" json:"default_withdrawlimit_daily"`
	DefaultMaxWithdraw        float64 `db:"default_max_withdraw" json:"default_max_withdraw"`

	BaseModel
}

func (this *Currency) GetTableName() string {
	return `currency`
}

func (this *Currency) GetByCurrencyChain(currency string, chain string) *Currency {
	result := Currency{}
	if notFound := go_mysql.MysqlHelper.SelectFirst(&result, this.GetTableName(), `*`, map[string]interface{}{
		`currency`:  currency,
		`chain`:     chain,
		`is_banned`: 0,
	}); notFound {
		return nil
	}
	return &result
}
