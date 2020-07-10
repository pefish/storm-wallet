package model

import (
	"github.com/pefish/go-mysql"
)

var DepositAddressModel = DepositAddress{}

type DepositAddress struct {
	UserId       uint64 `json:"user_id"`
	Address      string `json:"address"`
	Path         string `json:"path"`
	Tag          string `json:"tag"`
	Series       string `json:"series"`
	AddressIndex int64  `json:"address_index"`
	IsDeleted    int64  `json:"is_deleted"`
	BaseModel
}

func (this *DepositAddress) GetTableName() string {
	return `deposit_address`
}

func (this *DepositAddress) GetByUserIdSeriesIndex(userId uint64, series string, index uint64) *DepositAddress {
	result := DepositAddress{}
	if notFound := go_mysql.MysqlHelper.MustSelectFirstByStr(&result, this.GetTableName(), `*`, `
where user_id = ? and series = ? and is_deleted = 0 and address_index = ?
`, userId, series, index); notFound {
		return nil
	}
	return &result
}

func (this *DepositAddress) GetByUserIdSeriesAddress(userId uint64, series, address, tag string) *DepositAddress {
	result := DepositAddress{}
	if notFound := go_mysql.MysqlHelper.MustSelectFirstByStr(&result, this.GetTableName(), `*`, `
where user_id = ? and series = ? and is_deleted = 0 and address = ? and tag = ?
`, userId, series, address, tag); notFound {
		return nil
	}
	return &result
}

func (this *DepositAddress) Insert(userId uint64, address string, path string, series string, index uint64, tag string) {
	go_mysql.MysqlHelper.MustRawExec(
		`insert ignore into deposit_address (user_id, address, path, tag, series, address_index) values (?,?,?,?,?,?)`,
		userId,
		address,
		path,
		tag,
		series,
		index,
	)
}
