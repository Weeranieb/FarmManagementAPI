package model

import "github.com/shopspring/decimal"

type SellDetail struct {
	Id           int             `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	SellId       int             `json:"sellId" gorm:"column:sell_id"`
	Size         string          `json:"size" gorm:"column:size"`
	FishType     string          `json:"fishType" gorm:"column:fish_type"`
	Amount       decimal.Decimal `json:"amount" gorm:"column:amount"`
	FishUnit     string          `json:"fishUnit" gorm:"column:fish_unit"`
	PricePerUnit decimal.Decimal `json:"pricePerUnit" gorm:"column:price_per_unit"`
	BaseModel
}

func (SellDetail) TableName() string {
	return "sell_details"
}
