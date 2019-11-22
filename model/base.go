package model

type BaseModel struct {
	Id        uint64 `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
