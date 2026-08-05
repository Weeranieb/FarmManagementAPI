package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type FeedPriceHistory struct {
	Id               int             `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	FeedCollectionId int             `json:"feedCollectionId" gorm:"column:feed_collection_id"`
	Price            decimal.Decimal `json:"price" gorm:"column:price"`
	PriceUpdatedDate time.Time       `json:"priceUpdatedDate" gorm:"column:price_updated_date"`
	BaseModel
}
