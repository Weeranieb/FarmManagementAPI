package model

import "github.com/shopspring/decimal"

type FeedCollection struct {
	Id       int              `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ClientId int              `json:"clientId" gorm:"column:client_id"`
	Name     string           `json:"name" gorm:"column:name"`
	Unit     string           `json:"unit" gorm:"column:unit"`
	FeedType string           `json:"feedType" gorm:"column:feed_type;not null;default:pellet"`
	Fcr      *decimal.Decimal `json:"fcr,omitempty" gorm:"column:fcr"`
	BaseModel
}

type FeedCollectionPage struct {
	FeedCollection
	LatestPrice            decimal.Decimal `json:"latestPrice" gorm:"column:latest_price"`
	LatestPriceUpdatedDate *string         `json:"latestPriceUpdatedDate" gorm:"column:latest_price_updated_date"`
}
