package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type DailyLog struct {
	Id                int             `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ActivePondId      int             `json:"activePondId" gorm:"column:active_pond_id;not null"`
	FeedDate          time.Time       `json:"feedDate" gorm:"column:feed_date;type:date;not null"`
	FreshMorning      decimal.Decimal `json:"freshMorning" gorm:"column:fresh_morning;not null;default:0"`
	FreshEvening      decimal.Decimal `json:"freshEvening" gorm:"column:fresh_evening;not null;default:0"`
	PelletMorning     decimal.Decimal `json:"pelletMorning" gorm:"column:pellet_morning;not null;default:0"`
	PelletEvening     decimal.Decimal `json:"pelletEvening" gorm:"column:pellet_evening;not null;default:0"`
	DeathFishCount    int             `json:"deathFishCount" gorm:"column:death_fish_count;not null;default:0"`
	TouristCatchCount *int            `json:"touristCatchCount" gorm:"column:tourist_catch_count"`
	BaseModel
}

func (DailyLog) TableName() string {
	return "daily_logs"
}
