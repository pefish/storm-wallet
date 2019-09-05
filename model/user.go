package model

import (
	"github.com/pefish/go-mysql"
)

var UserModel = User{}

type User struct {
	Mobile           string `db:"mobile" json:"mobile"`
	IsBanned         int64  `db:"is_banned" json:"is_banned"`
	BaseModel
}

func (this *User) GetTableName() string {
	return `user`
}

func (this *User) GetByUserIdIsBanned(userId uint64, isBanned bool) *User {
	result := User{}
	if notFound := go_mysql.MysqlHelper.SelectFirst(&result, this.GetTableName(), `*`, map[string]interface{}{
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
