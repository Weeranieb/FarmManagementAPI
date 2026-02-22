package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Activity struct {
	Id             int             `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ActivePondId   int             `json:"activePondId" gorm:"column:active_pond_id"`
	ToActivePondId *int            `json:"toActivePondId,omitempty" gorm:"column:to_active_pond_id"`
	Mode           string          `json:"mode" gorm:"column:mode"`
	MerchantId     *int            `json:"merchantId,omitempty" gorm:"column:merchant_id"`
	Amount         int             `json:"amount" gorm:"column:amount;not null"`
	FishType       string          `json:"fishType" gorm:"column:fish_type;not null"`
	FishWeight     decimal.Decimal `json:"fishWeight" gorm:"column:fish_weight;not null"`
	FishUnit       string          `json:"fishUnit" gorm:"column:fish_unit;not null"`
	PricePerUnit   decimal.Decimal `json:"pricePerUnit" gorm:"column:price_per_unit;not null"`
	ActivityDate   time.Time       `json:"activityDate" gorm:"column:activity_date"`
	BaseModel
}

func (Activity) TableName() string {
	return "activities"
}
