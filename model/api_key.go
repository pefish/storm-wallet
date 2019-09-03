package model

import (
	"github.com/pefish/go-mysql"
)

var ApiKeyModel = ApiKey{}

type ApiKey struct {
	UserId    uint64 `db:"user_id" json:"user_id"`
	ApiKey    string `db:"api_key" json:"api_key"`
	ApiSecret string `db:"api_secret" json:"api_secret"`
	Ip        string `db:"ip" json:"ip"`
	Type      int64  `db:"type" json:"type"`
	IsDeleted int64  `db:"is_deleted" json:"is_deleted"`
	BaseModel
}

func (this *ApiKey) GetTableName() string {
	return `api_key`
}

func (this *ApiKey) GetByApiKey(apiKey string) *ApiKey {
	apiKeyModel := ApiKey{}
	if notFound := go_mysql.MysqlHelper.SelectFirst(&apiKeyModel, this.GetTableName(), `*`, map[string]interface{}{
		`api_key`:    apiKey,
		`is_deleted`: 0,
	}); notFound {
		return nil
	}
	return &apiKeyModel
}
