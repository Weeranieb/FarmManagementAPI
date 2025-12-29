package model

type Pond struct {
	Id     int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	FarmId int    `json:"farmId" gorm:"column:farm_id"`
	Code   string `json:"code" gorm:"column:code"`
	Name   string `json:"name" gorm:"column:name"`
	BaseModel
}

