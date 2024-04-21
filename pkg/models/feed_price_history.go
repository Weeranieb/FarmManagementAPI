package models

import (
	"errors"
	"time"
)

// FeedPriceHistory represents the history of feed prices.
type FeedPriceHistory struct {
	Id               int       `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	FeedCollectionId int       `json:"feedCollectionId" gorm:"column:FeedCollectionId"`
	Price            float64   `json:"price" gorm:"column:Price"`
	PriceUpdatedDate time.Time `json:"priceUpdatedDate" gorm:"column:PriceUpdatedDate"`
	Base
}

type AddFeedPriceHistory struct {
	FeedCollectionId int       `json:"feedCollectionId" gorm:"column:FeedCollectionId"`
	Price            float64   `json:"price" gorm:"column:Price"`
	PriceUpdatedDate time.Time `json:"priceUpdatedDate" gorm:"column:PriceUpdatedDate"`
}

// Validation Add
func (a AddFeedPriceHistory) Validation() error {
	if a.FeedCollectionId == 0 {
		return errors.New(ErrFeedCollectionIdEmpty)
	}
	if a.Price == 0 {
		return errors.New(ErrPriceEmpty)
	}
	if a.PriceUpdatedDate.IsZero() {
		return errors.New(ErrPriceUpdatedDateEmpty)
	}
	return nil
}

// Transfer Add
func (a AddFeedPriceHistory) Transfer(feedPriceHistory *FeedPriceHistory) {
	feedPriceHistory.FeedCollectionId = a.FeedCollectionId
	feedPriceHistory.Price = a.Price
	feedPriceHistory.PriceUpdatedDate = a.PriceUpdatedDate
}

const (
	ErrFeedCollectionIdEmpty = "feedCollectionId is empty"
	ErrPriceEmpty            = "price is empty"
	ErrPriceUpdatedDateEmpty = "priceUpdatedDate is empty"
)
