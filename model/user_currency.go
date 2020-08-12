package model

import (
	"fmt"
	"github.com/pefish/go-mysql"
	"github.com/pefish/go-reflect"
	"strings"
)

var UserCurrencyModel = UserCurrency{}

type UserCurrency struct {
	UserId             uint64  `json:"user_id"`
	CurrencyId         uint64  `json:"currency_id"`
	WithdrawLimitDaily float64 `json:"withdraw_limit_daily"`
	MaxWithdrawAmount  float64 `json:"max_withdraw_amount"`
	WithdrawCheckLimit float64 `json:"withdraw_check_limit"`
	IsDeleted          int64   `json:"is_deleted"`
	BaseModel
}

func (this *UserCurrency) GetTableName() string {
	return `user_currency`
}

func (this *UserCurrency) GetCurrencyOfUserByCurrencyChain(userId uint64, currency string, chain string) *Currency {
	result := Currency{}
	select_ := strings.Join(go_reflect.Reflect.GetValuesInTagFromStruct(&result, `json`), `,b.`)
	if notFound := go_mysql.MysqlInstance.MustRawSelectFirst(&result, fmt.Sprintf(`
select b.%s from user_currency a
left join currency b
on a.currency_id = b.id
where b.currency = ? and b.chain = ? and a.is_deleted = 0 and b.is_banned = 0 and a.user_id = ?
`, select_), currency, chain, userId); notFound {
		return nil
	}
	return &result
}

func (this *UserCurrency) GetByUserIdCurrencyId(userId uint64, currencyId uint64) *UserCurrency {
	result := UserCurrency{}
	if notFound := go_mysql.MysqlInstance.MustSelectFirst(&result, this.GetTableName(), `*`, map[string]interface{}{
		`user_id`:     userId,
		`currency_id`: currencyId,
		`is_deleted`:  0,
	}); notFound {
		return nil
	}
	return &result
}

func (this *UserCurrency) ListByUserId(results interface{}, userId uint64) {
	go_mysql.MysqlInstance.MustSelect(results, this.GetTableName(), `*`, map[string]interface{}{
		`user_id`:    userId,
		`is_deleted`: 0,
	})
}
