package model

import (
	"github.com/pefish/go-mysql"
)

var ResponseKeyModel = ResponseKey{}

type ResponseKey struct {
	UserId     uint64 `db:"user_id" json:"user_id"`
	PublicKey  string `db:"public_key" json:"public_key"`
	PrivateKey string `db:"private_key" json:"private_key"`
	IsDeleted  int64  `db:"is_deleted" json:"is_deleted"`
	BaseModel
}

func (this *ResponseKey) GetTableName() string {
	return `response_key`
}

func (this *ResponseKey) GetByUserId(userId uint64) *ResponseKey {
	responseKeyModel := ResponseKey{}
	if notFound := go_mysql.MysqlHelper.SelectFirst(&responseKeyModel, this.GetTableName(), `*`, map[string]interface{}{
		`user_id`:    userId,
		`is_deleted`: 0,
	}); notFound {
		return nil
	}
	return &responseKeyModel
}
