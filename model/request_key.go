package model

import (
	"github.com/pefish/go-mysql"
)

var RequestKeyModel = RequestKey{}

type RequestKey struct {
	UserId    uint64 `db:"user_id" json:"user_id"`
	PublicKey string `db:"public_key" json:"public_key"`
	Ip        string `db:"ip" json:"ip"`
	Type      int64  `db:"type" json:"type"`
	IsDeleted int64  `db:"is_deleted" json:"is_deleted"`
	BaseModel
}

func (this *RequestKey) GetTableName() string {
	return `request_key`
}

func (this *RequestKey) GetByPubKey(pubKey string) *RequestKey {
	requestKeyModel := RequestKey{}
	if notFound := go_mysql.MysqlHelper.SelectFirst(&requestKeyModel, this.GetTableName(), `*`, map[string]interface{}{
		`public_key`: pubKey,
		`is_deleted`: 0,
	}); notFound {
		return nil
	}
	return &requestKeyModel
}
