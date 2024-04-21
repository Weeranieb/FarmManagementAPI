package models

import (
	"errors"
	"time"
)

// DailyFeed represents a daily feed in the system.
type DailyFeed struct {
	Id               int       `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	ActivePondId     int       `json:"activePondId" gorm:"column:ActivePondId"`
	FeedCollectionId int       `json:"feedCollectionId" gorm:"column:FeedCollectionId"`
	Amount           float64   `json:"amount" gorm:"column:Amount"`
	FeedDate         time.Time `json:"feedDate" gorm:"column:FeedDate"`
	Base
}

type AddDailyFeed struct {
	ActivePondId     int       `json:"activePondId" gorm:"column:ActivePondId"`
	FeedCollectionId int       `json:"feedCollectionId" gorm:"column:FeedCollectionId"`
	Amount           float64   `json:"amount" gorm:"column:Amount"`
	FeedDate         time.Time `json:"feedDate" gorm:"column:FeedDate"`
}

// Validation Add
func (a AddDailyFeed) Validation() error {
	if a.ActivePondId == 0 {
		return errors.New(ErrActivePondIdEmpty)
	}
	if a.FeedCollectionId == 0 {
		return errors.New(ErrFeedCollectionIdEmpty)
	}
	if a.Amount == 0 {
		return errors.New(ErrAmountEmpty)
	}
	if a.FeedDate.IsZero() {
		return errors.New(ErrFeedDateEmpty)
	}
	return nil
}

// Transfer Add
func (a AddDailyFeed) Transfer(dailyFeed *DailyFeed) {
	dailyFeed.ActivePondId = a.ActivePondId
	dailyFeed.FeedCollectionId = a.FeedCollectionId
	dailyFeed.Amount = a.Amount
	dailyFeed.FeedDate = a.FeedDate
}

const (
	ErrFeedDateEmpty = "feedDate is empty"
)
