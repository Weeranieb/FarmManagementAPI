package models

import "errors"

// FeedCollection represents a feed collection in the system.
type FeedCollection struct {
	Id       int    `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	ClientId int    `json:"clientId" gorm:"column:ClientId"`
	Code     string `json:"code" gorm:"column:Code"`
	Name     string `json:"name" gorm:"column:Name"`
	Unit     string `json:"unit" gorm:"column:Unit"`
	Base
}

type AddFeedCollection struct {
	Code string `json:"code" gorm:"column:Code"`
	Name string `json:"name" gorm:"column:Name"`
	Unit string `json:"unit" gorm:"column:Unit"`
}
type CreateFeedRequest struct {
	Code             string                `json:"code" gorm:"column:Code"`
	Name             string                `json:"name" gorm:"column:Name"`
	Unit             string                `json:"unit" gorm:"column:Unit"`
	FeedPriceHistory []AddFeedPriceHistory `json:"feedPriceHistory"`
}

// Validation Add
func (a AddFeedCollection) Validation() error {
	if a.Code == "" {
		return errors.New(ErrCodeEmpty)
	}
	if a.Name == "" {
		return errors.New(ErrNameEmpty)
	}
	if a.Unit == "" {
		return errors.New(ErrUnitEmpty)
	}
	return nil
}

// Transfer Add
func (a AddFeedCollection) Transfer(feedCollection *FeedCollection) {
	feedCollection.Code = a.Code
	feedCollection.Name = a.Name
	feedCollection.Unit = a.Unit
}

const (
	ErrUnitEmpty = "unit is empty"
)
