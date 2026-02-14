package model

import "time"

type FeedPriceHistory struct {
	Id               int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	FeedCollectionId int       `json:"feedCollectionId" gorm:"column:feed_collection_id"`
	Price            float64   `json:"price" gorm:"column:price"`
	PriceUpdatedDate time.Time `json:"priceUpdatedDate" gorm:"column:price_updated_date"`
	BaseModel
}
