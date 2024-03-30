package models

import "time"

// Activity represents an activity.
type Activity struct {
	Id             int       `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	ActivePondId   int       `json:"activePondId" gorm:"column:ActivePondId"`
	ToActivePondId *int      `json:"toActivePondId" gorm:"column:ToActivePondId"`
	Mode           string    `json:"mode" gorm:"column:Mode"`
	MerchantId     int       `json:"merchantId" gorm:"column:MerchantId"`
	Amount         int       `json:"amount" gorm:"column:Amount"`
	FishType       string    `json:"fishType" gorm:"column:FishType"`
	FishWeight     float64   `json:"fishWeight" gorm:"column:FishWeight"`
	PricePerUnit   float64   `json:"pricePerUnit" gorm:"column:PricePerUnit"`
	ActivityDate   time.Time `json:"activityDate" gorm:"column:ActivityDate"`
	Base
}
