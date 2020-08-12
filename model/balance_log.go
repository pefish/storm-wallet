package model

import (
	"errors"
	"fmt"
	"github.com/pefish/go-error"
	"github.com/pefish/go-mysql"
)

var BalanceLogModel = BalanceLog{}

var (
	LogType_BalanceChange uint64 = 1
	LogType_FreezeChange  uint64 = 2

	Type_Withdraw uint64 = 1
	Type_Deposit  uint64 = 2
	Type_Air      uint64 = 3
)

type BalanceLog struct {
	UserId     uint64  `json:"user_id"`
	CurrencyId uint64  `json:"currency_id"`
	Amount     float64 `json:"amount"`
	Type       int64   `json:"type"`
	LogType    int64   `json:"log_type"`
	Ref_id     int64   `json:"ref_id"`

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
	if notFound := go_mysql.MysqlInstance.MustRawSelectFirst(&result, fmt.Sprintf(`
select sum(if(log_type = 1, amount, 0)) as avail, sum(if(log_type = 2, amount, 0)) as freeze
from %s where user_id = ? and currency_id = ?
`, this.GetTableName()), userId, currencyId); notFound {
		result.Avail = &zero
		result.Freeze = &zero
		return &result
	}
	if result.Avail == nil {
		result.Avail = &zero
	}
	if result.Freeze == nil {
		result.Freeze = &zero
	}
	return &result
}

func (this *BalanceLog) ListAllBalanceByUserIdForStruct(results interface{}, userId uint64) {
	go_mysql.MysqlInstance.MustRawSelect(results, fmt.Sprintf(`
select sum(if(a.log_type = 1, a.amount, 0)) as avail, sum(if(a.log_type = 2, a.amount, 0)) as freeze, b.currency, b.chain
from %s a
left join %s b
on a.currency_id = b.id
where a.user_id = ?
group by b.currency, b.chain
`, this.GetTableName(), CurrencyModel.GetTableName()), userId)
}

func (this *BalanceLog) Freeze(tran *go_mysql.MysqlClass, userId uint64, currencyId uint64, amount string, type_ uint64, refId uint64) {
	_, num := tran.MustInsert(this.GetTableName(), map[string]interface{}{
		`user_id`:     userId,
		`currency_id`: currencyId,
		`amount`:      amount,
		`type`:        type_,
		`log_type`:    LogType_FreezeChange,
		`ref_id`:      refId,
	})
	if num == 0 {
		go_error.ThrowInternal(errors.New(`freeze error`))
	}
}
