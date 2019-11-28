package model

import "github.com/pefish/go-mysql"

var CurrencyModel = Currency{}

type Currency struct {
	Currency                      string  `json:"currency"`
	Chain                         string  `json:"chain"`
	ContractAddress               string  `json:"contract_address"`
	Decimals                      int64   `json:"decimals"`
	ScanLimit                     float64 `json:"scan_limit"`
	ReservedValue                 float64 `json:"reserved_value"`
	DepositConfirmationThreshold  int64   `json:"deposit_confirmation_threshold"`
	WithdrawConfirmationThreshold int64   `json:"withdraw_confirmation_threshold"`
	Series                        string  `json:"series"`
	WithdrawInoutNum              int64   `json:"withdraw_inout_num"`
	NetworkFee                    float64 `json:"network_fee"`
	ScanInoutNum                  int64   `json:"scan_inout_num"`
	HasTag                        int64   `json:"has_tag"`
	IsWithdrawEnable              int64   `json:"is_withdraw_enable"`
	IsDepositEnable               int64   `json:"is_deposit_enable"`
	IsBanned                      int64   `json:"is_banned"`
	MaxTagLength                  int64   `json:"max_tag_length"`
	DefaultWithdrawLimitDaily     float64 `json:"default_withdrawlimit_daily"`
	DefaultMaxWithdraw            float64 `json:"default_max_withdraw"`

	BaseModel
}

func (this *Currency) GetTableName() string {
	return `currency`
}

func (this *Currency) GetByCurrencyChain(currency string, chain string) *Currency {
	result := Currency{}
	if notFound := go_mysql.MysqlHelper.SelectFirst(&result, this.GetTableName(), `id`, map[string]interface{}{
		`currency`:  currency,
		`chain`:     chain,
		`is_banned`: 0,
	}); notFound {
		return nil
	}
	return &result
}

func (this *Currency) GetIdByCurrencyChain(currency string, chain string, userId uint64) *Currency {
	result := Currency{}
	if notFound := go_mysql.MysqlHelper.SelectFirst(&result, this.GetTableName(), `*`, map[string]interface{}{
		`currency`: currency,
		`chain`:    chain,
		`user_id`:  userId,
	}); notFound {
		return nil
	}
	return &result
}
