package model

import (
	"fmt"
	"github.com/pefish/go-error"
	"github.com/pefish/go-mysql"
	"github.com/pefish/go-time"
)

var WithdrawModel = Withdraw{}

type Withdraw struct {
	RequestId     string  `json:"request_id"`
	UserId        uint64  `json:"user_id"`
	CurrencyId    uint64  `json:"currency_id"`
	Currency      string  `json:"currency"`
	Chain         string  `json:"chain"`
	Amount        float64 `json:"amount"`
	FromAddress   string  `json:"from_address"`
	ToAddress     string  `json:"to_address"`
	Status        int64   `json:"status"`
	Height        string  `json:"height"`
	BlockId       string  `json:"block_id"`
	TxId          string  `json:"tx_id"`
	TxHex         *string `json:"tx_hex"`
	Confirmations int64   `json:"confirmations"`
	OutputIndex   int64   `json:"output_index"`
	Nonce         int64   `json:"nonce"`
	NetworkFee    float64 `json:"network_fee"`
	Tag           string  `json:"tag"`
	Mark          *string `json:"mark"`

	BaseModel
}

func (this *Withdraw) GetTableName() string {
	return `withdraw`
}

func (this *Withdraw) ListByUserIdChainTxIdForStruct(results interface{}, userId uint64, chain string, txId string) {
	go_mysql.MysqlHelper.Select(results, this.GetTableName(), `*`, map[string]interface{}{
		`user_id`: userId,
		`chain`:   chain,
		`tx_id`:   txId,
	})
}

func (this *Withdraw) ListByUserIdTxIdForStruct(results interface{}, userId uint64, txId string) {
	go_mysql.MysqlHelper.Select(results, this.GetTableName(), `*`, map[string]interface{}{
		`user_id`: userId,
		`tx_id`:   txId,
	})
}

func (this *Withdraw) GetByUserIdRequestId(userId uint64, requestId string) *Withdraw {
	result := Withdraw{}
	if notFound := go_mysql.MysqlHelper.SelectFirst(&result, this.GetTableName(), `*`, map[string]interface{}{
		`user_id`:    userId,
		`request_id`: requestId,
	}); notFound {
		return nil
	}
	return &result
}

type WithdrewTotalStruct struct {
	Sum *string `json:"sum"`
}

func (this *Withdraw) GetWithdrewTotalOfToday(userId uint64, currency string, chain string) string {
	beginOfTodayTime := go_time.Time.GetFormatTimeFromTimeObj(go_time.Time.GetLocalBeginTimeOfToday(), `0000-00-00 00:00:00`)
	endOfTodayTime := go_time.Time.GetFormatTimeFromTimeObj(go_time.Time.GetLocalEndTimeOfToday(), `0000-00-00 00:00:00`)
	sumStruct := WithdrewTotalStruct{}
	go_mysql.MysqlHelper.RawSelectFirst(&sumStruct, fmt.Sprintf(`
select sum(amount) as sum 
from %s 
where 
	user_id = ?
	and currency = ?
	and chain = ?
	and created_at between ? and ?
`, this.GetTableName()), userId, currency, chain, beginOfTodayTime, endOfTodayTime)
	if sumStruct.Sum == nil {
		return `0`
	}
	return *sumStruct.Sum
}

func (this *Withdraw) Insert(tran go_mysql.MysqlClass, requestId string, userId uint64, currencyId uint64, currency string, chain string, amount string, status uint64, address string, memo string) uint64 {
	id, num := tran.Insert(this.GetTableName(), map[string]interface{}{
		`request_id`:  requestId,
		`user_id`:     userId,
		`currency_id`: currencyId,
		`currency`:    currency,
		`chain`:       chain,
		`amount`:      amount,
		`status`:      status,
		`to_address`:  address,
		`tag`:         memo,
	})
	if num == 0 {
		go_error.ThrowInternal(`insert error`)
	}
	return uint64(id)
}
