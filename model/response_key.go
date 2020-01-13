package model

import (
	"github.com/pefish/go-mysql"
)

var ResponseKeyModel = ResponseKey{}

type ResponseKey struct {
	UserId     uint64 `json:"user_id"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	IsDeleted  int64  `json:"is_deleted"`
	BaseModel
}

func (this *ResponseKey) GetTableName() string {
	return `response_key`
}

func (this *ResponseKey) GetByUserId(userId uint64) *ResponseKey {
	responseKeyModel := ResponseKey{}
	if notFound := go_mysql.MysqlHelper.MustSelectFirst(&responseKeyModel, this.GetTableName(), `*`, map[string]interface{}{
		`user_id`:    userId,
		`is_deleted`: 0,
	}); notFound {
		return nil
	}
	return &responseKeyModel
}
