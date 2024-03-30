package models

import "time"

// DailyFeed represents a daily feed in the system.
type DailyFeed struct {
	Id               int       `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	ActivePondId     int       `json:"activePondId" gorm:"column:ActivePondId"`
	FeedCollectionId int       `json:"feedCollectionId" gorm:"column:FeedCollectionId"`
	Amount           float64   `json:"amount" gorm:"column:Amount"`
	FeedDate         time.Time `json:"feedDate" gorm:"column:FeedDate"`
	Base
}
