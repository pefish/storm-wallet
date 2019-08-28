package model

var ApiKeyModel = ApiKey{}

type ApiKey struct {
	UserId    uint64  `db:"user_id" json:"user_id"`
	ApiKey    string  `db:"api_key" json:"api_key"`
	ApiSecret string  `db:"api_secret" json:"api_secret"`
	Ip        string  `db:"ip" json:"ip"`
	Type      int64   `db:"type" json:"type"`
	DeletedAt *string `db:"deleted_at" json:"deleted_at"`
	BaseModel
}

func (this *ApiKey) GetTableName() string {
	return `api_key`
}
