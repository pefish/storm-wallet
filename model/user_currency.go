package model

import (
	"fmt"
	"github.com/pefish/go-mysql"
	"github.com/pefish/go-reflect"
	"strings"
)

var UserCurrencyModel = UserCurrency{}

type UserCurrency struct {
	UserId             uint64  `db:"user_id" json:"user_id"`
	CurrencyId         uint64  `db:"currency_id" json:"currency_id"`
	WithdrawLimitDaily float64 `db:"withdraw_limit_daily" json:"withdraw_limit_daily"`
	MaxWithdrawAmount  float64 `db:"max_withdraw_amount" json:"max_withdraw_amount"`
	WithdrawCheckLimit float64 `db:"withdraw_check_limit" json:"withdraw_check_limit"`
	IsDeleted          int64   `db:"is_deleted" json:"is_deleted"`
	BaseModel
}

func (this *UserCurrency) GetTableName() string {
	return `user_currency`
}

func (this *UserCurrency) GetCurrencyOfUserByCurrencyChain(userId uint64, currency string, chain string) *Currency {
	result := Currency{}
	select_ := strings.Join(go_reflect.Reflect.GetValuesInTagFromStruct(&result, `db`), `,b.`)
	if notFound := go_mysql.MysqlHelper.RawSelectFirst(&result, fmt.Sprintf(`
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
	if notFound := go_mysql.MysqlHelper.SelectFirst(&result, this.GetTableName(), `*`, map[string]interface{}{
		`user_id`:     userId,
		`currency_id`: currencyId,
		`is_deleted`:  0,
	}); notFound {
		return nil
	}
	return &result
}
