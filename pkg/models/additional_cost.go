package models

import "errors"

// AdditionalCost represents an additional cost in the system.
type AdditionalCost struct {
	Id         int     `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	ActivityId int     `json:"activityId" gorm:"column:ActivityId"`
	Title      string  `json:"title" gorm:"column:Title"`
	Cost       float64 `json:"cost" gorm:"column:Cost"`
	Base
}

type AddAdditionalCostRequest struct {
	ActivityId int     `json:"activityId"`
	Title      string  `json:"title"`
	Cost       float64 `json:"cost"`
}

// Validation Add
func (a AddAdditionalCostRequest) Validation() error {
	if a.ActivityId == 0 {
		return errors.New(ErrActivityIdEmpty)
	}
	if a.Title == "" {
		return errors.New(ErrTitleEmpty)
	}
	if a.Cost == 0 {
		return errors.New(ErrCostEmpty)
	}
	return nil
}

// Transfer Add
func (a AddAdditionalCostRequest) Transfer(additionalCost *AdditionalCost) {
	additionalCost.ActivityId = a.ActivityId
	additionalCost.Title = a.Title
	additionalCost.Cost = a.Cost
}

const (
	ErrActivityIdEmpty = "activity id is empty"
	ErrTitleEmpty      = "title is empty"
	ErrCostEmpty       = "cost is empty"
)
