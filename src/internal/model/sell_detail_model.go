package model

import "github.com/shopspring/decimal"

type SellDetail struct {
	Id              int             `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	SellId          int             `json:"sellId" gorm:"column:sell_id"`
	FishSizeGradeId int             `json:"fishSizeGradeId" gorm:"column:fish_size_grade_id"`
	Weight          decimal.Decimal `json:"weight" gorm:"column:weight"`
	PricePerUnit    decimal.Decimal `json:"pricePerUnit" gorm:"column:price_per_unit"`
	FishCount       *int            `json:"fishCount,omitempty" gorm:"column:fish_count"`
	BaseModel
}

func (SellDetail) TableName() string {
	return "sell_details"
}
