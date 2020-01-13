package model

import (
	"github.com/pefish/go-mysql"
)

var TeamModel = Team{}

type Team struct {
	Mobile             string `json:"mobile"`
	WithdrawConfirmUrl string `json:"withdraw_confirm_url"`
	IsBanned           int64  `json:"is_banned"`
	BaseModel
}

func (this *Team) GetTableName() string {
	return `user`
}

func (this *Team) GetByUserIdIsBanned(userId uint64, isBanned bool) *Team {
	result := Team{}
	if notFound := go_mysql.MysqlHelper.MustSelectFirst(&result, this.GetTableName(), `*`, map[string]interface{}{
		`is_banned`: func(isBanned bool) int {
			if isBanned {
				return 1
			} else {
				return 0
			}
		}(isBanned),
		`id`: userId,
	}); notFound {
		return nil
	}
	return &result
}
