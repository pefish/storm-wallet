package model

import "github.com/pefish/go-mysql"

var WalletConfigModel = WalletConfig{}

type WalletConfig struct {
	Chain     string  `json:"chain"`
	Address   string  `json:"address"`
	Path      string  `json:"path"`
	MaxLimit  float64 `json:"max_limit"`
	MinLimit  float64 `json:"min_limit"`
	Type      int64   `json:"type"`
	IsDeleted string  `json:"is_deleted"`

	BaseModel
}

func (this *WalletConfig) GetTableName() string {
	return `wallet_config`
}

func (this *WalletConfig) GetByChainType(chain string, type_ uint64) *WalletConfig {
	result := WalletConfig{}
	if notFound := go_mysql.MysqlHelper.SelectFirst(&result, this.GetTableName(), `*`, map[string]interface{}{
		`type`:       type_,
		`chain`:      chain,
		`is_deleted`: 0,
	}); notFound {
		return nil
	}
	return &result
}
