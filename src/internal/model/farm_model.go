package model

type Farm struct {
	Id       int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ClientId int    `json:"clientId" gorm:"column:client_id"`
	Code     string `json:"code" gorm:"column:code"`
	Name     string `json:"name" gorm:"column:name"`
	BaseModel
}
