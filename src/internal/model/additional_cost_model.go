package model

import "github.com/shopspring/decimal"

type AdditionalCost struct {
	Id         int             `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ActivityId int             `json:"activityId" gorm:"column:activity_id"`
	Title      string          `json:"title" gorm:"column:title"`
	Cost       decimal.Decimal `json:"cost" gorm:"column:cost"`
	BaseModel
}

func (AdditionalCost) TableName() string {
	return "additional_costs"
}
