package models

import (
	"errors"
	"time"
)

// ActivePond represents a pond that is currently active.
type ActivePond struct {
	Id        int        `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	PondId    int        `json:"pondId" gorm:"column:PondId"`
	StartDate time.Time  `json:"startDate" gorm:"column:StartDate"`
	EndDate   *time.Time `json:"endDate" gorm:"column:EndDate"`
	IsActive  bool       `json:"isActive" gorm:"column:IsActive"`
	Base
}

type AddActivePond struct {
	PondId    int       `json:"pondId" gorm:"column:PondId"`
	StartDate time.Time `json:"startDate" gorm:"column:StartDate"`
}

// Validation Add
func (a AddActivePond) Validation() error {
	if a.PondId == 0 {
		return errors.New(ErrPondIdEmpty)
	}
	if a.StartDate.IsZero() {
		return errors.New(ErrStartDateEmpty)
	}
	return nil
}

// Transfer Add
func (a AddActivePond) Transfer(activePond *ActivePond) {
	activePond.PondId = a.PondId
	activePond.StartDate = a.StartDate
}

const (
	ErrPondIdEmpty    = "pond id is empty"
	ErrStartDateEmpty = "start date is empty"
)

type PondWithActive struct {
	Id           int    `json:"id" gorm:"column:Id"`
	Code         string `json:"code" gorm:"column:Code"`
	Name         string `json:"name" gorm:"column:Name"`
	ActivePondId *int   `json:"activePondId" gorm:"column:ActivePondId"`
	HasHistory   bool   `json:"hasHistory" gorm:"column:HasHistory"`
}
