package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type DailyFeed struct {
	Id               int             `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ActivePondId     int             `json:"activePondId" gorm:"column:active_pond_id;not null"`
	FeedCollectionId int             `json:"feedCollectionId" gorm:"column:feed_collection_id;not null"`
	FeedDate         time.Time       `json:"feedDate" gorm:"column:feed_date;type:date;not null"`
	MorningAmount    decimal.Decimal `json:"morningAmount" gorm:"column:morning_amount;not null;default:0"`
	EveningAmount    decimal.Decimal `json:"eveningAmount" gorm:"column:evening_amount;not null;default:0"`
	BaseModel
}

func (DailyFeed) TableName() string {
	return "daily_feeds"
}
