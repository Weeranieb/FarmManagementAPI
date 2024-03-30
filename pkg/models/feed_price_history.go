package models

import "time"

// FeedPriceHistory represents the history of feed prices.
type FeedPriceHistory struct {
	Id               int       `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	FeedCollectionId int       `json:"feedCollectionId" gorm:"column:FeedCollectionId"`
	Price            float64   `json:"price" gorm:"column:Price"`
	PriceUpdatedDate time.Time `json:"priceUpdatedDate" gorm:"column:PriceUpdatedDate"`
	Base
}
