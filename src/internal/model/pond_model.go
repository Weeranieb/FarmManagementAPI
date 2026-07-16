package model

import "github.com/shopspring/decimal"

type Pond struct {
	Id     int                 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	FarmId int                 `json:"farmId" gorm:"column:farm_id"`
	Name   string              `json:"name" gorm:"column:name"`
	Status string              `json:"status" gorm:"column:status;default:'maintenance'"`
	Area   decimal.NullDecimal `json:"area" gorm:"column:area_rai"`
	BaseModel
}

// PondCountPerClient holds the pond total grouped by client_id (via farms).
type PondCountPerClient struct {
	ClientId int   `gorm:"column:client_id"`
	Total    int64 `gorm:"column:total"`
}
