package model

import (
	"fmt"
	"github.com/pefish/go-error"
	"github.com/pefish/go-mysql"
)

var BalanceLogModel = BalanceLog{}

var (
	LogType_BalanceRealAdd uint64 = 1
	LogType_BalanceRealSub uint64 = 2
	LogType_FreezeAdd      uint64 = 3
	LogType_FreezeSub      uint64 = 4

	Type_Withdraw uint64 = 1
	Type_Deposit  uint64 = 2
	Type_Air      uint64 = 3
)

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
select sum(if(log_type in (1,2), amount, 0)) as avail, sum(if(log_type in (3,4), amount, 0)) as freeze
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

func (this *BalanceLog) Freeze(tran go_mysql.MysqlClass, userId uint64, currencyId uint64, amount string, type_ uint64, refId uint64) {
	_, num := tran.Insert(this.GetTableName(), map[string]interface{}{
		`user_id`:     userId,
		`currency_id`: currencyId,
		`amount`:      `-` + amount,
		`type`:        type_,
		`log_type`:    LogType_FreezeAdd,
		`ref_id`:      refId,
	})
	if num == 0 {
		go_error.ThrowInternal(`freeze error`)
	}
}
