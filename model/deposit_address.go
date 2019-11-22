package model

import (
	"github.com/pefish/go-mysql"
)

var DepositAddressModel = DepositAddress{}

type DepositAddress struct {
	UserId       uint64 `json:"user_id"`
	Address      string `json:"address"`
	Path         string `json:"path"`
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
	if notFound := go_mysql.MysqlHelper.SelectFirstByStr(&result, this.GetTableName(), `*`, `
where user_id = ? and series = ? and is_deleted = 0 and address_index = ?
`, userId, series, index); notFound {
		return nil
	}
	return &result
}

func (this *DepositAddress) GetByUserIdSeriesAddress(userId uint64, series string, address string) *DepositAddress {
	result := DepositAddress{}
	if notFound := go_mysql.MysqlHelper.SelectFirstByStr(&result, this.GetTableName(), `*`, `
where user_id = ? and series = ? and is_deleted = 0 and address = ?
`, userId, series, address); notFound {
		return nil
	}
	return &result
}

func (this *DepositAddress) Insert(userId uint64, address string, path string, series string, index uint64) {
	go_mysql.MysqlHelper.RawExec(
		`insert into deposit_address (user_id, address, path, series, address_index) values (?,?,?,?,?)`,
		userId,
		address,
		path,
		series,
		index,
	)
}
