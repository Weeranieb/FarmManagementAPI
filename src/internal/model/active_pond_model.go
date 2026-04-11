package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type ActivePond struct {
	Id          int             `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	PondId      int             `json:"pondId" gorm:"column:pond_id"`
	StartDate   time.Time       `json:"startDate" gorm:"column:start_date"`
	EndDate     *time.Time      `json:"endDate,omitempty" gorm:"column:end_date"`
	IsActive    bool            `json:"isActive" gorm:"column:is_active"`
	TotalCost   decimal.Decimal `json:"totalCost" gorm:"column:total_cost"`
	TotalProfit decimal.Decimal `json:"totalProfit" gorm:"column:total_profit"`
	NetResult   decimal.Decimal `json:"netResult" gorm:"column:net_result"`
	TotalFish   int             `json:"totalFish" gorm:"column:total_fish"`
	FishTypes   []string        `json:"fishTypes" gorm:"column:fish_types;serializer:json"`
	// Default feed collections for daily log UI when a month has no rows yet (at most one fresh + one pellet per active cycle).
	DefaultFreshFeedCollectionId  *int `json:"defaultFreshFeedCollectionId,omitempty" gorm:"column:default_fresh_feed_collection_id"`
	DefaultPelletFeedCollectionId *int `json:"defaultPelletFeedCollectionId,omitempty" gorm:"column:default_pellet_feed_collection_id"`
	BaseModel
}

func (ActivePond) TableName() string {
	return "active_ponds"
}
