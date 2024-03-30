package models

// SellDetail represents a sell detail in the system.
type SellDetail struct {
	Id           int     `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	SellId       int     `json:"sellId" gorm:"column:SellId"`
	Size         string  `json:"size" gorm:"column:Size"`
	FishType     string  `json:"fishType" gorm:"column:FishType"`
	Amount       float64 `json:"amount" gorm:"column:Amount"`
	FishUnit     string  `json:"fishUnit" gorm:"column:FishUnit"`
	PricePerUnit float64 `json:"pricePerUnit" gorm:"column:PricePerUnit"`
	Base
}
