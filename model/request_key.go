package model

import (
	"github.com/pefish/go-mysql"
)

var RequestKeyModel = RequestKey{}

type RequestKey struct {
	UserId    uint64 `json:"user_id"`
	PublicKey string `json:"public_key"`
	Ip        string `json:"ip"`
	Type      int64  `json:"type"`
	IsDeleted int64  `json:"is_deleted"`
	BaseModel
}

func (this *RequestKey) GetTableName() string {
	return `request_key`
}

func (this *RequestKey) GetByPubKey(pubKey string) *RequestKey {
	requestKeyModel := RequestKey{}
	if notFound := go_mysql.MysqlHelper.MustSelectFirst(&requestKeyModel, this.GetTableName(), `*`, map[string]interface{}{
		`public_key`: pubKey,
		`is_deleted`: 0,
	}); notFound {
		return nil
	}
	return &requestKeyModel
}
