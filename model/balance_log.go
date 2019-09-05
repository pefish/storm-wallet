package model

import (
	"fmt"
	"github.com/pefish/go-error"
	"github.com/pefish/go-mysql"
)

var BalanceLogModel = BalanceLog{}

type BalanceLog struct {
	UserId     uint64  `db:"user_id" json:"user_id"`
	CurrencyId uint64  `db:"currency_id" json:"currency_id"`
	Amount     float64 `db:"amount" json:"amount"`
	Type       int64   `db:"type" json:"type"`
	LogType    int64   `db:"log_type" json:"log_type"`
	Ref_id     int64   `db:"ref_id" json:"ref_id"`

	BaseModel
}

func (this *BalanceLog) GetTableName() string {
	return `balance_log`
}

type Balance struct {
	Avail  *string `json:"avail"`
	Freeze *string `json:"freeze"`
}

func (this *BalanceLog) GetBalanceByUserIdCurrencyId(userId uint64, currencyId uint64) *Balance {
	result := Balance{}
	zero := `0`
	if notFound := go_mysql.MysqlHelper.RawSelectFirst(&result, fmt.Sprintf(`
select sum(if(log_type = 2, 0, amount)) as avail, sum(if(log_type = 1, 0, amount)) as freeze
from %s where user_id = ? and currency_id = ?
`, this.GetTableName()), userId, currencyId); notFound {
		result.Avail = &zero
		result.Freeze = &zero
	}
	if result.Avail == nil {
		result.Avail = &zero
	}
	if result.Freeze == nil {
		result.Freeze = &zero
	}
	return &result
}

func (this *BalanceLog) Freeze(userId uint64, currencyId uint64, amount string, type_ uint64, refId uint64) {
	_, num := go_mysql.MysqlHelper.Insert(this.GetTableName(), map[string]interface{}{
		`user_id`:     userId,
		`currency_id`: currencyId,
		`amount`:      `-` + amount,
		`type`:        type_,
		`log_type`:    2,
		`ref_id`:      refId,
	})
	if num == 0 {
		go_error.ThrowInternal(`freeze error`)
	}
}
