package model

type Pond struct {
	Id     int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	FarmId int    `json:"farmId" gorm:"column:farm_id"`
	Name   string `json:"name" gorm:"column:name"`
	Status string `json:"status" gorm:"column:status;default:'maintenance'"`
	BaseModel
}
